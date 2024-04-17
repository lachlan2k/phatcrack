package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/roles"
	"github.com/lachlan2k/phatcrack/api/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/apitypes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const nonceCookieName = "_oauth_state_nonce"

type oidcUtils struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
	provider *oidc.Provider
}

type oauthState struct {
	Nonce    string
	Redirect string
}

func getOidcUtils(ctx context.Context) (*oidcUtils, error) {
	ac := config.Get().Auth

	provider, err := oidc.NewProvider(ctx, ac.OIDC.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %v", err)
	}

	config := &oauth2.Config{
		ClientID:     ac.OIDC.ClientID,
		ClientSecret: ac.OIDC.ClientSecret,
		RedirectURL:  ac.OIDC.RedirectURL,

		Endpoint: provider.Endpoint(),
		Scopes:   append([]string{oidc.ScopeOpenID, "email", "profile"}, ac.OIDC.AdditionalScopes...),
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	return &oidcUtils{
		config:   config,
		provider: provider,
		verifier: verifier,
	}, nil
}

func handleOIDCStart(sessHandler auth.SessionHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isOIDCAuthAllowed() {
			return echo.NewHTTPError(http.StatusBadRequest, "OIDC auth is not enabled")
		}

		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
		utils, err := getOidcUtils(ctx)

		if err != nil {
			return util.GenericServerError(err)
		}

		nonceBuff := make([]byte, 16)
		_, err = rand.Read(nonceBuff)
		if err != nil {
			return util.GenericServerError(fmt.Errorf("failed to generated nonce: %v", err))
		}

		nonceStr := base64.RawURLEncoding.EncodeToString(nonceBuff)
		state := oauthState{
			Nonce:    nonceStr,
			Redirect: c.QueryParam("redir"),
		}

		stateBuff, err := json.Marshal(state)
		if err != nil {
			return util.GenericServerError(fmt.Errorf("failed to marshal state for oauth: %v", err))
		}

		stateStr := string(stateBuff)

		c.SetCookie(&http.Cookie{
			Name:     nonceCookieName,
			Value:    nonceStr,
			Expires:  time.Now().Add(5 * time.Minute),
			Secure:   c.Scheme() == "https",
			Path:     "/",
			HttpOnly: true,
		})

		return c.Redirect(http.StatusFound, utils.config.AuthCodeURL(stateStr))
	}
}

func handleOIDCCallback(sessHandler auth.SessionHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isOIDCAuthAllowed() {
			return echo.NewHTTPError(http.StatusBadRequest, "OIDC auth is not enabled")
		}

		ac := config.Get().Auth

		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
		utils, err := getOidcUtils(ctx)
		if err != nil {
			return util.GenericServerError(err)
		}

		var state oauthState
		err = json.Unmarshal([]byte(c.QueryParam("state")), &state)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid state")
		}

		cookieNonce, err := c.Cookie(nonceCookieName)
		if err != nil || cookieNonce.Value == "" {
			return c.String(http.StatusBadRequest, "State cookie wasn't found: request likely expired")
		}

		if cookieNonce.Value != state.Nonce {
			return c.String(http.StatusBadRequest, "State nonce mismatch")
		}

		code := c.QueryParam("code")
		if code == "" {
			return c.String(http.StatusBadRequest, "No code was provided")
		}

		token, err := utils.config.Exchange(ctx, code)
		if err != nil {
			return util.GenericServerError(fmt.Errorf("failed perform oauth2 exchange: provided code was likely invalid: %v", err))
		}

		rawToken, ok := token.Extra("id_token").(string)
		if !ok {
			log.WithField("access_token", token.AccessToken).Debug("failed casting id_token")
			return util.GenericServerError(fmt.Errorf("couldn't cast id_token to string err: %v", err))
		}

		idToken, err := utils.verifier.Verify(ctx, rawToken)
		if err != nil {
			return util.GenericServerError(fmt.Errorf("couldn't verify token: %v", err))
		}

		claims := map[string]any{}
		err = idToken.Claims(&claims)
		if err != nil {
			return util.GenericServerError(fmt.Errorf("failed to map token claims: %v", err))
		}

		username, extractedRoles, err := extractUsernameAndRoles(claims, ac)
		if err != nil {
			return util.GenericServerError(err)
		}

		if len(ac.OIDC.RequiredRole) > 0 && !slices.Contains(extractedRoles, ac.OIDC.RequiredRole) {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this application")
		}

		// User is now verified, time to sign in
		existingUser, err := db.GetUserByUsername(username)

		var userToAuth *db.User = existingUser

		if err != nil && err != db.ErrNotFound {
			return util.GenericServerError(err)
		}

		if err == db.ErrNotFound || existingUser == nil {
			if !ac.OIDC.AutomaticUserCreation {
				AuditLog(c, log.Fields{
					"username": username,
				}, "User tried to sign in with OIDC, but they do not have an account, and automatic user creation is disabled")
				return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this application")
			}

			// Create a new user
			newUser, err := db.RegisterUserWithoutPassword(username, []string{roles.RoleStandard})
			if err != nil {
				return util.GenericServerError(fmt.Errorf("failed to create user: %v", err))
			}
			userToAuth = newUser
		}

		sessHandler.Start(c, auth.SessionData{
			UserID:          userToAuth.ID.String(),
			HasCompletedMFA: false,
		})

		response := apitypes.AuthLoginResponseDTO{
			User: apitypes.AuthCurrentUserDTO{
				ID:       userToAuth.ID.String(),
				Username: userToAuth.Username,
				Roles:    userToAuth.Roles,
			},
			IsAwaitingMFA:          userToAuth.HasRole(roles.RoleMFAEnrolled),
			RequiresPasswordChange: userToAuth.HasRole(roles.RoleRequiresPasswordChange),
			RequiresMFAEnrollment:  config.Get().Auth.IsMFARequired && !userToAuth.HasRole(roles.RoleMFAEnrolled) && !userToAuth.HasRole(roles.RoleMFAExempt),
		}

		logMessage := "Session started"
		if response.IsAwaitingMFA {
			logMessage += ", user needs to complete MFA"
		}
		if response.RequiresPasswordChange {
			logMessage += ", user needs to change password"
		}
		if response.RequiresMFAEnrollment {
			logMessage += ", user needs to enroll MFA"
		}

		AuditLog(c, log.Fields{
			"authenticated_username": userToAuth.Username,
			"auth_type":              "oidc",
		}, logMessage)

		return c.JSON(http.StatusOK, response)
	}
}

func extractUsernameAndRoles(claims map[string]any, ac config.AuthConfig) (string, []string, error) {
	claimUsernameAny, ok := claims[ac.OIDC.UsernameClaim]
	if !ok {
		return "", nil, fmt.Errorf("username claim %q was empty on sso", ac.OIDC.UsernameClaim)
	}

	claimUsernameStr, ok := claimUsernameAny.(string)
	if !ok {
		return "", nil, fmt.Errorf("couldn't cast username claim %q to string", ac.OIDC.UsernameClaim)
	}
	if len(claimUsernameStr) == 0 {
		return "", nil, fmt.Errorf("username claim %q was empty", ac.OIDC.UsernameClaim)
	}

	roleClaim := claims[ac.OIDC.RolesClaim]

	rolesAny, ok := roleClaim.([]any)
	if !ok {
		return "", nil, fmt.Errorf("couldn't cast roles %v (type of %T) to []any", roleClaim, roleClaim)
	}

	roles := make([]string, len(rolesAny))
	for i, v := range rolesAny {
		roleStr, ok := v.(string)
		if !ok {
			return "", nil, fmt.Errorf("failed to cast role item [%d] %v (type of %T) to string", i, v, v)
		}
		roles[i] = roleStr
	}

	return claimUsernameStr, roles, nil
}
