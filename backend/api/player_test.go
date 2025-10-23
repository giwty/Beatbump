package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)



func TestPlayer(t *testing.T) {

	//
	// Setup
	e := echo.New()
	//req := httptest.NewRequest(http.MethodGet, "/", nil)

	//c.SetPath("/users/:email")
	q := make(url.Values)
	q.Set("videoId", "uJdu4Lfy8aI")
	//q.Set("playlistId", "")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, PlayerEndpointHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s := rec.Body.String()
		fmt.Print(s)
		//assert.Equal(t, userJSON, rec.Body.String())
	}

}

func TestSearch(t *testing.T) {
	e := echo.New()
	//req := httptest.NewRequest(http.MethodGet, "/", nil)

	//c.SetPath("/users/:email")
	q := make(url.Values)
	q.Set("q", "depeche mode")
	//q.Set("playlistId", "")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, SearchEndpointHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		s := rec.Body.String()
		fmt.Print(s)
		//assert.Equal(t, userJSON, rec.Body.String())
	}
}
