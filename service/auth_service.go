package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type callbackResult struct {
	SessionID string
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
	InvalidateSession(ctx context.Context, sessionId string) error
}

type auth struct {
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	authConfig   *config.AuthConfig
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

	// Get the auth config
	authConfig := c.Authentication

	return &auth{
		provider: provider,
		verifier: verifier,
		oauth2Config: &oauth2.Config{
			ClientID:     authConfig.ClientID,
			ClientSecret: authConfig.ClientSecret,
			RedirectURL:  authConfig.RedirectURL,
			// Discovery returns the OAuth2 endpoints.
			Endpoint: provider.Endpoint(),
			Scopes:   authConfig.Scopes,
		},
		db:         db,
		sess:       sess,
		authConfig: &authConfig,
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
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	rawIDToken, ok := accessToken.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("No id_token field in oauth2 token.")
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify ID Token: %w", err)
	}

	// Parse the standard claims
	var idTokenClaims claims
	if err := idToken.Claims(&idTokenClaims); err != nil {
		return nil, fmt.Errorf("Failed to parse ID Token claims: %w", err)
	}

	// Parse the roles claim
	var allClaims map[string]interface{}
	if err := idToken.Claims(&allClaims); err != nil {
		return nil, fmt.Errorf("Failed to parse ID Token claims: %w", err)
	}
	if roles, ok := allClaims["roles"].([]string); ok {
		idTokenClaims.Roles = roles
	}

	// Check if "admin" role is present
	isAdmin := false
	for _, role := range idTokenClaims.Roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	// Create user in DBFS if not exists
	var user *dbfs.User
	err = ctxutil.GetTransaction(ctx, a.db).Transaction(func(tx *gorm.DB) error {
		if user, err = dbfs.GetUserByMappedId(tx, idTokenClaims.Subject); err != nil {
			// User doesn't exist, create
			if user, err = dbfs.CreateNewUser(tx,
				idTokenClaims.Email, idTokenClaims.Name, dbfs.AccountTypeEndUser,
				idTokenClaims.Subject, "TODO", accessToken.AccessToken, rawIDToken, "server",
			); err != nil {
				return fmt.Errorf("Failed to create new user: %w", err)
			}
		}

		// Fetch and reassign user's groups
		groups, err := dbfs.GetGroupsByRoleNames(tx, idTokenClaims.Roles)
		if err != nil {
			return fmt.Errorf("Failed to fetch groups: %w", err)
		}

		if err := user.SetGroups(tx, groups); err != nil {
			return fmt.Errorf("Failed to assign groups: %w", err)
		}

		// Set admin bit
		if isAdmin {
			user.AccountType = dbfs.AccountTypeAdmin
		} else {
			user.AccountType = dbfs.AccountTypeEndUser
		}
		if err := tx.Save(&user).Error; err != nil {
			return fmt.Errorf("Failed to assign user type: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Create session ID from user ID
	sessionId, err := a.sess.Create(ctx, user.UserId, time.Hour*24*7)
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %w", err)
	}

	// Return session ID
	return &callbackResult{
		SessionID: sessionId,
	}, nil
}

func (a *auth) InvalidateSession(ctx context.Context, sessionId string) error {
	return a.sess.Invalidate(ctx, sessionId)
}
