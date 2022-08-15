package web_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/OhanaFS/ohana/web"
	"github.com/stretchr/testify/assert"
)

func TestWebEmbed(t *testing.T) {
	assert := assert.New(t)

	handler, err := web.GetHandler()
	assert.NoError(err)

	// Load index page
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Result().StatusCode)
	assert.Contains(w.Result().Header.Get("Content-Type"), "text/html")
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(err)
	assert.NotEmpty(body)

	// Find the script file
	re := regexp.MustCompile(`(?m)(/assets/.*?.js)`)
	scriptPath := re.FindString(string(body))
	assert.NotEmpty(scriptPath, "Could not find bundled javascript")

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", scriptPath, nil)
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Result().StatusCode)
	assert.Contains(w.Result().Header.Get("Content-Type"), "javascript")
	jsbody, err := io.ReadAll(w.Result().Body)
	assert.NoError(err)
	assert.NotEmpty(jsbody)

	assert.NotEqual(body, jsbody)
}
