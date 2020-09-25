package holberton

import (
	"holberton/api/app/models"
	"holberton/api/logger"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func (h *Holberton) project(id string) (*models.Project, error) {
	var err error

	visited := false

	project := &models.Project{
		ID: id,
	}

	localURL := "/project/" + id
	projectURL := "https://intranet.hbtn.io/projects/" + id
	_, err = h.page.Goto(projectURL)
	if err != nil {
		logger.Log2(err, "could not goto")
		return nil, err
	}

	html, _ := h.page.Content()
	h.setHtml(html, localURL)

	h.collector.OnHTML("article", func(article *colly.HTMLElement) {
		if visited {
			return
		}

		categoryP := article.DOM.Find("body > main > article > p.sm-gap")
		categoryTitle := strings.Trim(categoryP.Text(), "\t\n ")
		project.Category = categoryTitle

		projectTitle := article.DOM.Find("h1.gap")
		project.Title = projectTitle.Text()

		scoreSelector := "body > main > article > div.gap.clean.well > ul > li:nth-child(2) >" +
			" ul > li:nth-child(3) > strong"
		projectScore := article.DOM.Find(scoreSelector)
		project.Score = projectScore.Text()

		tasksContainerSelector := "body > main > article > section"
		tasksContainer := article.DOM.Find(tasksContainerSelector)
		tasksContainer.Children().Each(func(i int, task *goquery.Selection) {
			if task.Is("div.quiz_question_item_container") {
				return
			}

			value, _ := task.Attr("data-role")
			taskID := value[4:]

			h4 := task.Find("h4.task")
			span := task.Find("h4 > span")

			title := strings.Replace(h4.Text(), span.Text(), "", 1)
			title = strings.Trim(title, "\t\n ")
			class := strings.Trim(span.Text(), "\t\n ")

			if class[0] == '#' { //Ex, #advanced to advanced
				class = class[1:]
			}

			project.Tasks = append(project.Tasks, models.Task{
				ID: taskID,
				Title: title,
				Type: class,
			})
		})

		visited = true
	})

	h.collector.Visit(h.ts.URL + localURL)

	return project, nil
}
