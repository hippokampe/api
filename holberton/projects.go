package holberton

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/utils"
	"github.com/mxschmitt/playwright-go"
	"github.com/pkg/errors"
)

func (hbtn *Holberton) GetProjects(email string) (models.Projects, error) {
	scope := "holberton"
	ctx, err := hbtn.getSession(email)
	if err != nil {
		return models.Projects{}, errors.Wrap(err, scope)
	}

	projects, err := hbtn.getProjects(*ctx.BrowserContext)
	if err != nil {
		return models.Projects{}, errors.Wrap(err, scope)
	}

	return projects, nil
}

func (hbtn *Holberton) getProjects(browserCtx playwright.BrowserContext) (models.Projects, error) {
	scope := "projects"
	page := browserCtx.Pages()[0]

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

func (hbtn *Holberton) extractProject(browserCtx playwright.BrowserContext, projectUrl string) (models.Project, error) {
	scope := "fill project"
	page := browserCtx.Pages()[0]

	_, err := page.Goto(projectUrl)
	if err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	html, err := page.Content()
	if err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	url := hbtn.setHtml(html)
	hbtn.collector.OnHTML("body > main:nth-child(3) > article:nth-child(3)", func(article *colly.HTMLElement) {
		titleSelector := "body > main:nth-child(3) > article:nth-child(3) > h1"
		categorySelector := "p.sm-gap:nth-child(4) > small:nth-child(1)"
		weightSelector := "body > main:nth-child(3) > article:nth-child(3) > p:nth-child(6) > em:nth-child(1) > small:nth-child(1)"

		title := article.DOM.Find(titleSelector)
		category := article.DOM.Find(categorySelector)
		weight := article.DOM.Find(weightSelector)

		fmt.Println(title.Text())
		fmt.Println(category.Text())
		fmt.Println(weight.Text())

	})

	err = hbtn.collector.Visit(url)
	if err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	return models.Project{}, nil
}
