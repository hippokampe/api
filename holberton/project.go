package holberton

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/utils"
	"github.com/mxschmitt/playwright-go"
	"github.com/pkg/errors"
)

func (hbtn *Holberton) getProject(browserCtx playwright.BrowserContext, id string) (models.Project, error) {
	scope := "project"
	page := browserCtx.Pages()[0]

	projectUrl := fmt.Sprintf("https://intranet.hbtn.io/projects/%s", id)
	if _, err := page.Goto(projectUrl); err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	html, err := page.Content()
	if err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	var project models.Project

	url := hbtn.setHtml(html)
	hbtn.collector.OnHTML("main > article", func(article *colly.HTMLElement) {
		titleSelector := "h1.gap"
		categorySelector := "main > article > p.sm-gap"
		scoreSelector := "div.gap.clean.well > ul > li:nth-child(2) > ul > li:nth-child(3) > strong"

		titleElement := article.DOM.Find(titleSelector)
		categoryElement := article.DOM.Find(categorySelector)
		scoreElement := article.DOM.Find(scoreSelector)

		// Calculating project score
		projectScore := utils.CleanString(scoreElement.Text())
		if projectScore == "" {
			projectScore = "0%"
		}

		project.ID = id
		project.Title = utils.CleanString(titleElement.Text())
		project.Category = utils.CleanString(categoryElement.Text())
		project.Score = projectScore

		// Iterating over the tasks
		article.ForEach("section.formatted-content > div", func(_ int, taskElement *colly.HTMLElement) {
			isTask := !strings.HasPrefix(taskElement.Attr("data-role"), "quiz_question")
			if isTask {
				task := hbtn.getTask(taskElement)
				project.Tasks = append(project.Tasks, task)
			}
		})
	})

	if err := hbtn.collector.Visit(url); err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	return project, nil
}
