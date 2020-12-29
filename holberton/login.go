package holberton

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/hippokampe/api/models"
	"github.com/mxschmitt/playwright-go"
	"github.com/pkg/errors"
)

func (hbtn *Holberton) login(browserCtx playwright.BrowserContext, credentials models.Login) (models.User, error) {
	scope := "login"
	page, err := browserCtx.NewPage()
	if err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	_, err = page.Goto("https://intranet.hbtn.io/auth/sign_in")
	if err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	if err := page.Fill("#user_login", credentials.Email); err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	if err := page.Fill("#user_password", credentials.Password); err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	if err := page.Click("#new_user > div.actions > input", playwright.PageClickOptions{
		Timeout: playwright.Int(1000 * 5),
	}); err != nil {
		selectorMsg := `waiting for selector "#new_user > div.actions > input"`

		if !strings.Contains(err.Error(), selectorMsg) {
			return models.User{}, errors.Wrap(err, scope)
		}

		return models.User{}, errors.Wrap(ErrBadCredentials, scope)
	}

	_, err = page.Goto("https://intranet.hbtn.io/users/my_profile")
	if err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	html, err := page.Content()
	if err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	url := hbtn.setHtml(html)
	user, err := hbtn.fillUser(url)
	if err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	user.Email = credentials.Email
	return user, nil
}

func (hbtn *Holberton) fillUser(url string) (models.User, error) {
	var user models.User

	hbtn.collector.OnHTML("body", func(body *colly.HTMLElement) {
		picSelector := "div.row:nth-child(4) > div:nth-child(1) > img:nth-child(1)"
		firstNameSelector := "#user_first_name"
		lastNameSelector := "#user_last_name"
		usernameSelector := "#user_preferred_name"
		idSelector := ".list-group-item > p:nth-child(3) > code:nth-child(1)"
		citySelector := "#user_address_city"

		pic := body.DOM.Find(picSelector)
		firstName := body.DOM.Find(firstNameSelector)
		lastName := body.DOM.Find(lastNameSelector)
		username := body.DOM.Find(usernameSelector)
		id := body.DOM.Find(idSelector)
		city := body.DOM.Find(citySelector)

		user.ProfileURL, _ = pic.Attr("src")
		user.FirstName, _ = firstName.Attr("value")
		user.LastName, _ = lastName.Attr("value")
		user.UserName, _ = username.Attr("value")
		user.ID = id.Text()
		user.City, _ = city.Attr("value")
	})

	if err := hbtn.collector.Visit(url); err != nil {
		return models.User{}, ErrServeFile
	}

	return user, nil
}
