package holberton

import (
	"github.com/gocolly/colly"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/logger"
	"github.com/mxschmitt/playwright-go"
)

func (h *Holberton) login(email, password string) (*models.User, error) {
	var err error
	user := &models.User{
		Email:    email,
		Password: password,
	}

	_, err = h.page.Goto(BaseUrl + "/auth/sign_in")
	if err != nil {
		logger.Log2(err, "could not goto")
		return nil, err
	}

	err = h.page.Fill("#user_login", email)
	if err != nil {
		logger.Log2(err, "could not set the user[login]")
		return nil, err
	}

	err = h.page.Fill("#user_password", password)
	if err != nil {
		logger.Log2(err, "could not set the user[password]")
		return nil, err
	}

	err = h.page.Click("#new_user > div.actions > input")
	if err != nil {
		logger.Log2(err, "could not sent the information")
		return nil, err
	}

	_, err = h.page.Goto("https://intranet.hbtn.io/users/my_profile")

	exists, err := h.userExists(h.page, user)
	if err != nil {
		logger.Log2(err, "cannot check if user exists")
		return nil, err
	}

	if exists {
		return user, err
	}

	return nil, nil
}

func (h *Holberton) userExists(page *playwright.Page, user *models.User) (bool, error) {
	exists := false
	html, _ := page.Content()
	url := h.setHtml(html, "/login")

	selector := "#user_preferred_name"
	h.collector.OnHTML(selector, func(div *colly.HTMLElement) {
		user.Username = cleanString(div.Attr("value"))

		h.InternalStatus.Logged = true
		exists = true
	})

	h.collector.Visit(url)

	return exists, nil
}
