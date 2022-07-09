package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/OhanaFS/ohana/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type callbackResult struct {
	AccessToken string  `json:"access_token"`
	IdToken     string  `json:"id_token"`
	UserInfo    *claims `json:"user_info"`
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
}

type auth struct {
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
}

var state = "somestate"

func NewAuth(c *config.Config) (Auth, error) {
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
			Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
		},
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

	return &callbackResult{
		AccessToken: accessToken.AccessToken,
		IdToken:     rawIDToken,
		UserInfo:    &idTokenClaims,
	}, nil
}
