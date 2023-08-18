package auth

import (
	log "github.com/sirupsen/logrus"

	"github.com/NHAS/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/db"
)

type SessionHandler interface {
	CreateMiddleware() echo.MiddlewareFunc

	Start(echo.Context, SessionData) error
	Destroy(echo.Context) error
	Refresh(echo.Context) error
	Rotate(echo.Context) error // Rotate session cookie, to mitigate session fixation
	UpdateSessionData(c echo.Context, updateFunc func(*SessionData) error) error

	LogoutAllSessionsForUser(id string) error

	shouldSkip(echo.Context) bool
}

type SessionData struct {
	UserID          string
	HasCompletedMFA bool

	WebAuthnSession     *webauthn.SessionData
	PendingWebAuthnUser *webauthnUser
}

const sessionContextKey = "sess-data"
const sessionUserContextKey = "sess-user"

func SessionDataFromReq(c echo.Context) *SessionData {
	data, ok := c.Get(sessionContextKey).(SessionData)
	if !ok {
		return nil
	}

	return &data
}

func UserAndSessFromReq(c echo.Context) (*db.User, *SessionData) {
	sess := SessionDataFromReq(c)
	if sess == nil {
		return nil, nil
	}

	// Incase we've already retrived it in this context, don't bother fetching again
	existingUser, ok := c.Get(sessionUserContextKey).(*db.User)
	if ok && existingUser != nil {
		return existingUser, sess
	}

	user, err := db.GetUserByID(sess.UserID)
	if err != nil || user == nil {
		log.WithError(err).WithField("user_id", sess.UserID).Error("Failed to retrieve user's information from DB for session")
		return nil, nil
	}

	c.Set(sessionUserContextKey, user)
	return user, sess
}

func UserFromReq(c echo.Context) *db.User {
	u, _ := UserAndSessFromReq(c)
	return u
}
