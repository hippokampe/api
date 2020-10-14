package holberton

import (
	"github.com/hippokampe/configuration/v2/configuration"
	"net/http"
	"net/http/httptest"

	"github.com/hippokampe/api/logger"

	"github.com/gocolly/colly"

	"github.com/mxschmitt/playwright-go"
)

const (
	FIREFOX  = "firefox"
	CHROMIUN = "chromiun"
	WEBKIT   = "webkit"
)

const (
	BaseUrl = "https://intranet.hbtn.io"
)

type Holberton struct {
	pw             *playwright.Playwright
	browser        *playwright.Browser
	ts             *httptest.Server
	mux            *http.ServeMux
	collector      *colly.Collector
	page           *playwright.Page
	InternalStatus status
	Configuration  *configuration.InternalSettings
}

type status struct {
	Logged      bool
	VisitedURLS map[string]bool
	Started     bool
	Username    string
}

func NewSession(browserName string, config *configuration.InternalSettings) (*Holberton, error) {
	var err error

	pathBrowser, err := config.GetPathBrowser()
	if err != nil {
		return nil, err
	}

	holberton := &Holberton{
		Configuration: config,
	}

	browserOptions := playwright.BrowserTypeLaunchOptions{
		Headless:       playwright.Bool(true),
		ExecutablePath: playwright.String(pathBrowser),
	}

	holberton.pw, err = playwright.Run()
	if err != nil {
		return nil, err
	}

	// Selecting browser for scrapping
	switch browserName {
	case FIREFOX:
		holberton.browser, err = holberton.pw.Firefox.Launch(browserOptions)
	case CHROMIUN:
		holberton.browser, err = holberton.pw.Chromium.Launch(browserOptions)
	case WEBKIT:
		holberton.browser, err = holberton.pw.WebKit.Launch(browserOptions)
	default:
		return nil, logger.New("browser not available")
	}

	if err != nil {
		return nil, err
	}

	return holberton, nil
}

func (h *Holberton) CloseSession() {
	var err error
	if h.browser != nil {
		if err = h.browser.Close(); err != nil {
			logger.Log2(err, "could not close browser")
		}
	}

	if h.pw != nil {
		if err = h.pw.Stop(); err != nil {
			logger.Log2(err, "could not stop Playwright")
		}
	}
}
