package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/NHAS/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
)

type webauthnUser struct {
	Username      string                          `json:"username"`
	WebauthnID    []byte                          `json:"webauthn_id"`
	CredentialMap map[string]*webauthn.Credential `json:"credential_map"`
}

func (w webauthnUser) WebAuthnID() []byte {
	return w.WebauthnID
}

func (w webauthnUser) WebAuthnName() string {
	return w.Username
}

func (w webauthnUser) WebAuthnDisplayName() string {
	return w.Username
}

func (w webauthnUser) WebAuthnCredentials() []*webauthn.Credential {
	creds := make([]*webauthn.Credential, 0, len(w.CredentialMap))
	for _, cred := range w.CredentialMap {
		creds = append(creds, cred)
	}
	return creds
}

func (w webauthnUser) WebAuthnCredential(ID []byte) *webauthn.Credential {
	return w.CredentialMap[hex.EncodeToString(ID)]
}

func (w webauthnUser) WebAuthnIcon() string {
	return ""
}

func (w webauthnUser) PushCredential(c *webauthn.Credential) {
	w.CredentialMap[hex.EncodeToString(c.ID)] = c
}

var webauthnHandler webauthn.WebAuthn

func InitWebAuthn(baseURL url.URL) {
	webauthnHandler.Config = &webauthn.Config{
		RPDisplayName: "Phatcrack",
		RPID:          baseURL.Hostname(),
		RPOrigins:     []string{baseURL.String()},
	}
}

func MFAWebAuthnBeginRegister(c echo.Context, sessHandler SessionHandler) (marshalledJSONResponse []byte, userPresentableErr error, internalErr error) {
	user, sessData := UserAndSessFromReq(c)
	if user == nil {
		internalErr = errors.New("failed to get user from req")
		return
	}

	if user.HasRole(RoleMFAEnrolled) {
		userPresentableErr = errors.New("user already enrolled in MFA")
		return
	}

	if sessData.PendingWebAuthnUser != nil || sessData.WebAuthnSession != nil {
		userPresentableErr = errors.New("MFA registration already in progress")
		return
	}

	wID := make([]byte, 64)
	_, err := rand.Read(wID)
	if err != nil {
		internalErr = fmt.Errorf("couldn't generated user id for webauthn: %v", err)
		return
	}

	wUser := &webauthnUser{
		Username:      user.Username,
		WebauthnID:    wID,
		CredentialMap: make(map[string]*webauthn.Credential),
	}

	creation, webauthnSession, err := webauthnHandler.BeginRegistration(wUser)
	if err != nil {
		internalErr = fmt.Errorf("couldn't begin registration: %v", err)
		return
	}

	marshalled, err := json.Marshal(creation)
	if err != nil {
		internalErr = fmt.Errorf("couldn't marshal webauthn creation data: %v", err)
		return
	}

	err = sessHandler.UpdateSessionData(c, func(sd *SessionData) error {
		sd.WebAuthnSession = webauthnSession
		sd.PendingWebAuthnUser = wUser
		return nil
	})
	if err != nil {
		internalErr = fmt.Errorf("couldn't save webauthn session data in session: %v", err)
		return
	}

	return marshalled, nil, nil
}

func MFAWebAuthnFinishRegister(c echo.Context, sessHandler SessionHandler) (userPresentableErr error, internalErr error) {
	user, sessData := UserAndSessFromReq(c)
	if user == nil {
		internalErr = fmt.Errorf("failed to get user from req")
		return
	}

	if user.HasRole(RoleMFAEnrolled) {
		userPresentableErr = errors.New("user already enrolled in MFA")
		return
	}

	if sessData.PendingWebAuthnUser == nil || sessData.WebAuthnSession == nil {
		userPresentableErr = errors.New("registration process not started")
		return
	}

	credential, err := webauthnHandler.FinishRegistration(*sessData.PendingWebAuthnUser, *sessData.WebAuthnSession, c.Request())
	if err != nil {
		internalErr = fmt.Errorf("failed to finish webauthn registration: %v", err)
		return
	}

	sessData.PendingWebAuthnUser.PushCredential(credential)

	marshalledBytes, err := json.Marshal(sessData.PendingWebAuthnUser)
	if err != nil {
		internalErr = fmt.Errorf("failed to marshal webauthn user: %v", err)
		return
	}

	err = sessHandler.UpdateSessionData(c, func(sd *SessionData) error {
		sd.PendingWebAuthnUser = nil
		sd.WebAuthnSession = nil
		return nil
	})
	if err != nil {
		internalErr = fmt.Errorf("couldn't remove webauthn session data: %v", err)
		return
	}

	user.MFAType = MFATypeWebAuthn
	user.MFAData = marshalledBytes
	user.Roles = append(user.Roles, RoleMFAEnrolled)

	err = db.GetInstance().Save(user).Error
	if err != nil {
		internalErr = fmt.Errorf("failed to save user in database with new MFA data: %v", err)
		return
	}

	return nil, nil
}

func MFAWebAuthnBeginLogin(c echo.Context, sessHandler SessionHandler) (marshalledJSONResponse []byte, userPresentableErr error, internalErr error) {
	user := UserFromReq(c)
	if user == nil {
		internalErr = fmt.Errorf("failed to get user from req")
		return
	}

	if !user.HasRole(RoleMFAEnrolled) {
		userPresentableErr = fmt.Errorf("user is not enrolled in MFA")
		return
	}

	var wUser = &webauthnUser{}
	err := json.Unmarshal(user.MFAData, wUser)
	if err != nil {
		internalErr = fmt.Errorf("couldn't unmarshal user's MFA data: %v", err)
		return
	}

	credentialAssertion, webauthnSession, err := webauthnHandler.BeginLogin(wUser)
	if err != nil {
		internalErr = fmt.Errorf("failed to begin webauthn login: %v", err)
		return
	}

	err = sessHandler.UpdateSessionData(c, func(sd *SessionData) error {
		sd.WebAuthnSession = webauthnSession
		sd.PendingWebAuthnUser = wUser
		return nil
	})
	if err != nil {
		internalErr = fmt.Errorf("couldn't save webauthn session data in session: %v", err)
		return
	}

	marshalledJSONResponse, err = json.Marshal(credentialAssertion)
	if err != nil {
		internalErr = fmt.Errorf("failed to marshal credential assertion: %v", err)
		return
	}

	return
}

func MFAWebAuthnFinishLogin(c echo.Context, sessHandler SessionHandler) (userPresentableErr error, internalErr error) {
	user, sessData := UserAndSessFromReq(c)
	if user == nil {
		internalErr = fmt.Errorf("failed to get user from req")
		return
	}

	if sessData.PendingWebAuthnUser == nil || sessData.WebAuthnSession == nil {
		userPresentableErr = fmt.Errorf("MFA verification process not started")
		return
	}

	credential, err := webauthnHandler.FinishLogin(*sessData.PendingWebAuthnUser, *sessData.WebAuthnSession, c.Request())
	if err != nil {
		userPresentableErr = fmt.Errorf("MFA verification failed")
		c.Logger().Printf("User %s (%s) failed MFA verification with error: %v", user.Username, user.ID.String(), err)
		return
	}

	if credential.Authenticator.CloneWarning {
		userPresentableErr = fmt.Errorf("MFA verification failed (potential counter re-use)")
		c.Logger().Printf("User %s (%s) failed MFA verification due to clone warning", user.Username, user.ID.String())
		return
	}

	return
}
