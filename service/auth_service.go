package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OhanaFS/ohana/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

type clientRoles struct {
	Roles   string `json:"roles,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Scope   string `json:"scope,omitempty"`
	Fetched bool   `json:"fetched,omitempty"`
}

type Auth interface {
	SendRequest(ctx context.Context, rawAccessToken string) (string, error)
	Callback(ctx context.Context, code string, checkState string) (clientRoles, error)
}

type auth struct {
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
}

var state = "somestate"

func NewAuth(c *config.Config) (Auth, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, c.Authentication.ConfigURL)
	if err != nil {
		return nil, err
	}

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

func (a *auth) Callback(ctx context.Context, code string, checkState string) (clientRoles, error) {
	if checkState != state {
		return clientRoles{}, fmt.Errorf("state is invalid")
	}

	oauth2Token, err := a.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return clientRoles{}, err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return clientRoles{}, fmt.Errorf("No id_token field in oauth2 token.")
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return clientRoles{}, err
	}

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		return clientRoles{}, err
	}

	// Will remove this if not required later
	data, err := json.MarshalIndent(resp, "", "    ")
	data = data // avoid unused variable error

	if err != nil {
		return clientRoles{}, err
	}

	// service
	tokenString := oauth2Token.AccessToken
	roles, err := GetRolesFromJWT(tokenString)
	if err != nil {
		return clientRoles{}, err
	}
	return roles, err
}

// function to get the roles from the token
func GetRolesFromJWT(accesTokenString string) (clientRoles, error) {
	// Parsing the token just to extract values without validation.(as validated before)
	token, _, err := new(jwt.Parser).ParseUnverified(accesTokenString, jwt.MapClaims{})
	if err != nil {
		return clientRoles{}, err
	}

	userInfo := clientRoles{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		resourceAccess := claims["resource_access"]
		// TODO: doing casts like this is dangerous, maybe find a better way to do it - need to update
		resourceAccessClient := resourceAccess.(map[string]interface{})["DemoServiceClient"]
		resourceAccessRole := resourceAccessClient.(map[string]interface{})["roles"]
		for _, role := range resourceAccessRole.([]interface{}) {
			userInfo.Roles = userInfo.Roles + " " + role.(string)
		}
	} else {
		fmt.Println(err)
	}
	// extracting necessary info from the token
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	name := claims["name"].(string)
	email := claims["email"].(string)
	scope := claims["scope"].(string)

	userInfo.UserID = userID
	userInfo.Name = name
	userInfo.Email = email
	userInfo.Scope = scope
	userInfo.Fetched = true

	fmt.Print("USER SHIT", userInfo)

	return userInfo, nil
}
