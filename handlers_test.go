package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// "application"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_ping(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/ping", nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ping(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
