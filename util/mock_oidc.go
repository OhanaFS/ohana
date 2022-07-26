package mock_oidc

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

var (
	message string
	address string
	code    []string
)

// codeValid returns true if the code is valid
func codeValid(codeInReq string) bool {
	for _, v := range code {
		fmt.Println("found", v, "looking for", codeInReq)
		if v == codeInReq {
			return true
		}
	}
	return false
}

// invalidateCode removes the code from the list of valid codes
func invalidateCode(codeInReq string) {
	// not thread-safe
	var newListOfValidCodes []string
	for _, v := range code {
		if v != codeInReq {
			newListOfValidCodes = append(newListOfValidCodes, v)
		}
	}
	code = newListOfValidCodes
}

// handleAuth handles the /auth endpoint
func handleAuth(w http.ResponseWriter, r *http.Request) {

	whereToRedirect := r.URL.Query().Get("redirect_uri")
	if whereToRedirect == "" {
		whereToRedirect = "http://example.org?no=1234"
	}
	redirectURL, _ := url.Parse(whereToRedirect)
	params := redirectURL.Query()
	params.Set("state", r.URL.Query().Get("state"))

	generatedCode := fmt.Sprintf("%d", rand.Int())
	code = append(code, generatedCode)
	params.Set("code", generatedCode)

	redirectURL.RawQuery = params.Encode()

	w.Header().Set("Location", redirectURL.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// handleToken handles the /token endpoint
func handleToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	exchangeCode := r.FormValue("code")
	if !codeValid(exchangeCode) {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		invalidateCode(r.URL.Query().Get("code"))
		w.Write([]byte("{\"token\":\"a34a5f6\"}"))
	}
}

func handleUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Let the mock return whatever username you wish to
	username := os.Getenv("username")
	if username == "" {
		username = "johndoe@gmail.com"
	}
	w.Write([]byte(fmt.Sprintf("{\"name\":\"johndoe\",\"email\":\"a@a.com\",\"preferred_username\":\"%s\",\"sub\":\"62ccaea02\"}", username)))
}
