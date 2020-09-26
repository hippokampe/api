package holberton

import (
	"holberton/api/app/models"
	"holberton/api/logger"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func (h *Holberton) projects() (models.Projects, error) {
	var err error
	var projects []models.Project

	_, err = h.page.Goto("https://intranet.hbtn.io/dashboards/my_current_projects")
	if err != nil {
		logger.Log2(err, "could not goto")
		return nil, err
	}

	html, _ := h.page.Content()
	url := h.setHtml(html, "/projects")

	h.collector.OnHTML("body > main > article", func(element *colly.HTMLElement) {
		selector := "div.panel.panel-default"
		element.ForEach(selector, func(_ int, e *colly.HTMLElement) {
			h4 := e.DOM.Find("h4.panel-title")
			title := strings.Trim(h4.Text(), "\n\t ")

			projectsList := e.DOM.Find("li.list-group-item")
			projectsList.Each(func(_ int, project *goquery.Selection) {
				name := project.Find("a")       // Project title
				code := project.Find("code")    // Project code or id
				score := project.Find("strong") // Score of the project

				projects = append(projects, models.Project{
					Title:    name.Text(),
					ID:       code.Text(),
					Score:    score.Text(),
					Category: title,
				})
			})
		})
	})

	h.collector.Visit(url)

	return projects, nil
}
