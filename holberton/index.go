package holberton

import (
	"holberton/api/logger"
	"net/http"
	"net/http/httptest"

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
}

type status struct {
	Logged      bool
	VisitedURLS map[string]bool
	Started     bool
}

func NewSession(browserName string) (*Holberton, error) {
	var err error

	holberton := &Holberton{}

	browserOptions := playwright.BrowserTypeLaunchOptions{
		Headless:       playwright.Bool(true),
		ExecutablePath: playwright.String("/home/davixcky/.config/hbtn/firefox-1173/firefox/firefox"),
	}

	holberton.pw, err = playwright.Run()
	if err != nil {
		logger.Log2(err, "could not start playwright")
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
		logger.Log2(err, "could not launch browser")
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
