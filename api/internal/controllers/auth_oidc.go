package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"github.com/lachlan2k/phatcrack/api/internal/auth"
	"github.com/lachlan2k/phatcrack/api/internal/config"
	"github.com/lachlan2k/phatcrack/api/internal/util"
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

	provider, err := oidc.NewProvider(ctx, ac.OIDCIssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create oidc provider: %v", err)
	}

	config := &oauth2.Config{
		ClientID:     ac.OIDCClientID,
		ClientSecret: ac.OIDCClientSecret,
		RedirectURL:  ac.OIDCRedirectURL,

		Endpoint: provider.Endpoint(),
		Scopes:   append([]string{oidc.ScopeOpenID, "email", "profile"}, ac.OIDCAdditionalScopes...),
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

		return nil
	}
}
