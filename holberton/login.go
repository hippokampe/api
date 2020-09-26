package holberton

import (
	"holberton/api/app/models"
	"holberton/api/logger"
	"strings"

	"github.com/gocolly/colly"
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

	h.userExists(h.page, user)

	return user, nil
}

func (h *Holberton) userExists(page *playwright.Page, user *models.User) (bool, error) {
	exists := false
	html, _ := page.Content()
	url := h.setHtml(html, "/login")

	selector := "#user_preferred_name"
	h.collector.OnHTML(selector, func(div *colly.HTMLElement) {
		user.Username = strings.Trim(div.Attr("value"), "\t\n ")
		h.InternalStatus.Logged = true
		exists = true
	})

	h.collector.Visit(url)

	if !exists {
		return exists, logger.New("bad credentials")
	}

	return true, nil
}
