package holberton

import (
	"net/http"
	"net/http/httptest"

	"github.com/hippokampe/api/app/models"

	"github.com/gocolly/colly"
	"github.com/mxschmitt/playwright-go"
	"github.com/pkg/errors"
)

var (
	ErrBadCredentials   = errors.New("bad credentials")
	ErrServeFile        = errors.New("invalid access from colly to html content")
	ErrSessionNotExists = errors.New("session not found")
)

type Holberton struct {
	pw        *playwright.Playwright
	browser   playwright.Browser
	ts        *httptest.Server
	mux       *http.ServeMux
	sessions  map[string]*holbertonSession
	collector *colly.Collector
}

type holbertonSession struct {
	User           *models.User
	BrowserContext *playwright.BrowserContext
}
