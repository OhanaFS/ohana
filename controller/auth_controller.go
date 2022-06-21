package middleware

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/julienschmidt/httprouter"
)

type Res401Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"401"`
	Message  string `json:"message" example:"authorisation failed"`
}

//claims component of jwt contains many fields , we need only roles of OhanaFS
//"OhanaFS":{"OhanaFS":{"roles":["pets-admin","pet-details","pets-search"]}},
type Claims struct {
	ResourceAccess client `json:"resource_access,omitempty"`
	JTI            string `json:"jti,omitempty"`
}

type client struct {
	OhanaFS clientRoles `json:"OhanaFS,omitempty"`
}

type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

var RealmConfigURL string = "http://localhost:8080/auth/realms/TestingRealm"
var clientID string = "OhanaFS"

func IsAuthorizedJWT(h httprouter.Handle, role string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		rawAccessToken := r.Header.Get("Authorization") //for example oidc

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}
		ctx := oidc.ClientContext(context.Background(), client)
		provider, err := oidc.NewProvider(ctx, RealmConfigURL)
		if err != nil {
			authorisationFailed("authorisation failed while getting the provider: "+err.Error(), w, r)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: clientID,
		}
		verifier := provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(ctx, rawAccessToken)

		if err != nil {
			authorisationFailed("authorisation failed while verifying the token: "+err.Error(), w, r)
			return
		}

		var IDTokenClaims Claims // ID Token payload is just JSON.
		if err := idToken.Claims(&IDTokenClaims); err != nil {
			authorisationFailed("claims : "+err.Error(), w, r)
			return
		}

		//checking the roles
		user_access_roles := IDTokenClaims.ResourceAccess.OhanaFS.Roles
		fmt.Println("roles: ", user_access_roles)
		for _, b := range user_access_roles {
			if b == role {
				h(w, r, ps)
				return
			}
		}

		authorisationFailed("user not allowed to access this api", w, r)
	}
}

// Displaying if authorization failed
func authorisationFailed(message string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	data := Res401Struct{
		Status:   "FAILED",
		HTTPCode: http.StatusUnauthorized,
		Message:  message,
	}
	res, _ := json.Marshal(data)
	w.Write(res)
}
