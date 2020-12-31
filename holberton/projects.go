package holberton

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/hippokampe/api/models"
	"github.com/hippokampe/api/utils"
	"github.com/mxschmitt/playwright-go"
	"github.com/pkg/errors"
)

func (hbtn *Holberton) getProjects(browserCtx playwright.BrowserContext) (models.Projects, error) {
	scope := "projects"
	page, err := browserCtx.NewPage()
	if err != nil {
		return nil, errors.Wrap(err, scope)
	}

	defer page.Close()

	if _, err := page.Goto("https://intranet.hbtn.io/dashboards/my_current_projects"); err != nil {
		return models.Projects{}, errors.Wrap(err, scope)
	}

	html, err := page.Content()
	if err != nil {
		return models.Projects{}, err
	}

	var projects models.Projects

	url := hbtn.setHtml(html)
	hbtn.collector.OnHTML("main > article", func(articles *colly.HTMLElement) {
		selector := "div.panel.panel-default"
		articles.ForEach(selector, func(_ int, projectGroupsElement *colly.HTMLElement) {
			h4 := projectGroupsElement.DOM.Find("h4.panel-title")
			title := utils.CleanString(h4.Text())

			projectsList := projectGroupsElement.DOM.Find("li.list-group-item")
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

	if err := hbtn.collector.Visit(url); err != nil {
		return models.Projects{}, errors.Wrap(err, scope)
	}

	return projects, nil
}
