package holberton

import (
	"fmt"
	"holberton/api/app/models"
	"holberton/api/logger"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
)

func (h *Holberton) checkTask(id, taskId string) (*models.Task, error) {
	var err error
	var taskResult *models.Task
	taskResult = new(models.Task)

	visited := false
	localURL := "/project/" + id + "/task"
	projectURL := "https://intranet.hbtn.io/projects/" + id
	_, err = h.page.Goto(projectURL)
	if err != nil {
		logger.Log2(err, "could not goto")
		return nil, err
	}

	html, _ := h.page.Content()
	url := h.setHtml(html, localURL)

	h.collector.OnHTML("article", func(article *colly.HTMLElement) {
		if visited {
			return
		}

		taskSelector := fmt.Sprintf("#task-%s", taskId)
		task := article.DOM.Find("div" + taskSelector)

		h4, span := searchTitleTask(task)
		title, class := parseTitleTask(h4, span)
		taskResult.ID = taskId
		taskResult.Title = title
		taskResult.Type = class

		button := taskSelector +
			" > div.student_correction_requests > button.task_correction_modal.btn.btn-default"
		h.page.Click(button)

		selector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div > div"+
			" > div.modal-body > div.actions > center > input", taskId)
		h.page.Click(selector)

		resultSelector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div"+
			" > div > div.modal-body > div.result", taskId)
		h.page.WaitForSelector(resultSelector)

		html2, _ := h.page.Content()
		url2 := h.setHtml(html2, localURL+"/"+taskId)
		taskResult.ResultChecker = h.checker(url2, taskId)
		visited = true
	})

	h.collector.Visit(url)

	return taskResult, nil
}

func (h *Holberton) checker(url, taskId string) *models.ResultChecker {
	resultsChecker := &models.ResultChecker{}

	c2 := colly.NewCollector()

	selector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div > div > div.modal-body > div.result", taskId)
	h.page.Click(selector)

	c2.OnHTML(selector, func(element *colly.HTMLElement) {
		var checks []models.Check


		h4 := element.DOM.Find("h4")
		fmt.Println(h4.Text())

		selector := "div.check-inline"
		divs := element.DOM.Find(selector)
		divs.Each(func(i int, check *goquery.Selection) {
			checkID, _ := check.Attr("id")
			title, _ := check.Attr("title")
			typeCheck, status := parseCheck(title)

			checks = append(checks, models.Check{
				ID: checkID,
				Type: typeCheck,
				Status: status,
				Position: i,
			})
		})

		resultsChecker.Checks = checks
	})

	c2.Visit(url)
	return resultsChecker
}
