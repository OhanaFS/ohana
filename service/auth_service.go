package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"gorm.io/gorm"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type callbackResult struct {
	AccessToken  string        `json:"access_token"`
	IdToken      string        `json:"id_token"`
	RefreshToken string        `json:"refresh_token"`
	UserInfo     *claims       `json:"user_info"`
	TTL          time.Duration `json:"ttl"`
	SessionId    string        `json:"session_id"`
}

type claims struct {
	IssuedAt int64  `json:"iat,omitempty"`
	Expires  int64  `json:"exp,omitempty"`
	Subject  string `json:"sub,omitempty"`

	Roles []string `json:"roles,omitempty"`
	Name  string   `json:"name,omitempty"`
	Email string   `json:"email,omitempty"`
	Scope string   `json:"scope,omitempty"`
}

type Auth interface {
	SendRequest(ctx context.Context, rawAccessToken string) (string, error)
	Callback(ctx context.Context, code string, checkState string) (*callbackResult, error)
	InvalidateSession(ctx context.Context, userId string) error
}

type auth struct {
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
	db           *gorm.DB
	sess         Session
}

var state = "somestate"

func NewAuth(c *config.Config, db *gorm.DB, sess Session) (Auth, error) {
	// Fetch the provider configuration from the discovery endpoint.
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, c.Authentication.ConfigURL)
	if err != nil {
		return nil, err
	}

	// Configure the OAuth2 client.
	oidcConfig := &oidc.Config{
		ClientID: c.Authentication.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	return &auth{
		provider: provider,
		verifier: verifier,
		oauth2Config: &oauth2.Config{
			ClientID:     c.Authentication.ClientID,
			ClientSecret: c.Authentication.ClientSecret,
			RedirectURL:  c.Authentication.RedirectURL,
			// Discovery returns the OAuth2 endpoints.
			Endpoint: provider.Endpoint(),
			// "openid" is a required scope for OpenID Connect flows.
			Scopes: []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile", "email"},
		},
		db:   db,
		sess: sess,
	}, nil
}

func (a *auth) SendRequest(ctx context.Context, rawAccessToken string) (string, error) {
	if rawAccessToken == "" {
		return a.oauth2Config.AuthCodeURL(state), nil
	}

	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("rawAccessToken is invalid")
	}

	_, err := a.verifier.Verify(ctx, parts[1])

	if err != nil {
		return a.oauth2Config.AuthCodeURL(state), nil
	}
	return "Hello Ohanians", nil
}

func (a *auth) Callback(ctx context.Context, code string, checkState string) (*callbackResult, error) {
	if checkState != state {
		return nil, fmt.Errorf("state is invalid")
	}

	// Exchange the authorization code for an access token.
	accessToken, err := a.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	rawIDToken, ok := accessToken.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("No id_token field in oauth2 token.")
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify ID Token: %v", err)
	}

	var idTokenClaims claims
	if err := idToken.Claims(&idTokenClaims); err != nil {
		return nil, fmt.Errorf("Failed to parse ID Token claims: %v", err)
	}

	// refresh token

	refreshToken, ok := accessToken.Extra("refresh_token").(string)
	if !ok {
		return nil, fmt.Errorf("Failed to get refresh token: %v", err)
	}

	// account type for sending to createnewuser
	// need to include roles in the jwt
	var accountType int8 = 0
	/*if idTokenClaims.Roles[0] == "admin" {
		accountType = 1
	} else {
		accountType = 0
	}*/

	ttl := time.Duration(time.Hour * 24 * 7)
	uid := idTokenClaims.Subject
	tx := ctxutil.GetTransaction(ctx, a.db)
	var user *dbfs.User

	// validating user
	user, err = a.validate(ctx, uid) // if doesn't exist create a new user
	if err != nil {
		// create new user
		user, err = dbfs.CreateNewUser(tx, idTokenClaims.Email, idTokenClaims.Name, accountType, uid,
			refreshToken, accessToken.AccessToken, rawIDToken)
		if err != nil {
			return nil, fmt.Errorf("Failed to create new user: %v", err)
		}
	}
	// get id from user
	dbId := user.UserId

	fmt.Print("Hello world: ", dbId)
	// creating a new session
	sessionId, err := a.sess.Create(ctx, dbId, ttl)
	fmt.Print("sessionId: ", sessionId)
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %v", err)
	}

	return &callbackResult{
		AccessToken:  accessToken.AccessToken,
		IdToken:      rawIDToken,
		RefreshToken: refreshToken,
		UserInfo:     &idTokenClaims,
		SessionId:    sessionId,
		TTL:          ttl,
	}, nil
}

func (a *auth) validate(ctx context.Context, userId string) (*dbfs.User, error) {
	getGorm := ctxutil.GetTransaction(ctx, a.db)
	user, err := dbfs.GetUserById(getGorm, userId)
	if err.Error() == "user not found" {
		return nil, err
	}
	return user, nil
}

func (a *auth) InvalidateSession(ctx context.Context, sessionId string) error {
	err := a.sess.Invalidate(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("Failed to invalidate user: %v", err)
	}
	return nil
}
