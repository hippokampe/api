package holberton

import (
	"net/http"
	"net/http/httptest"

	"github.com/gocolly/colly"
	"github.com/mxschmitt/playwright-go"
)

func New() (*Holberton, error) {
	var err error

	hbtn := &Holberton{}

	hbtn.pw, err = playwright.Run()
	if err != nil {
		return nil, err
	}

	browserOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	}

	hbtn.browser, err = hbtn.pw.Firefox.Launch(browserOptions)
	if err != nil {
		return nil, err
	}

	hbtn.mux = http.NewServeMux()
	hbtn.ts = httptest.NewServer(hbtn.mux)
	hbtn.collector = colly.NewCollector()
	hbtn.sessions = make(map[string]*holbertonSession)

	return hbtn, nil
}

func (hbtn *Holberton) Close() error {
	var err error
	if hbtn.browser != nil {
		if err = hbtn.browser.Close(); err != nil {
			return err
		}
	}

	if hbtn.pw != nil {
		if err = hbtn.pw.Stop(); err != nil {
			return err
		}
	}

	return nil
}
