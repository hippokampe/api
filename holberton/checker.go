package holberton

import (
	"fmt"
	"holberton/api/logger"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
)

func (h *Holberton) checkTask(id, taskId string) error {
	var err error

	visited := false
	localURL := "/project/" + id + "/task"
	projectURL := "https://intranet.hbtn.io/projects/" + id
	_, err = h.page.Goto(projectURL)
	if err != nil {
		logger.Log2(err, "could not goto")
		return err
	}

	html, _ := h.page.Content()
	h.setHtml(html, localURL)

	h.collector.OnHTML("article", func(article *colly.HTMLElement) {
		if visited {
			return
		}

		taskSelector := fmt.Sprintf("#task-%s", taskId)
		task := article.DOM.Find("div" + taskSelector)
		title := task.Find("h4.task").Text()
		fmt.Println(title)
		//#task-1776 > div.student_correction_requests > button.task_correction_modal.btn.btn-default
		button := taskSelector +
			" > div.student_correction_requests > button.task_correction_modal.btn.btn-default"

		h.page.Click(button)
		selector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div > div >"+
			"div.modal-body > div.actions > center > input", taskId)
		h.page.Click(selector)
		resultSelector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div > div > div.modal-body > div.result", taskId)
		h.page.WaitForSelector(resultSelector)
		html2, _ := h.page.Content()
		h.setHtml(html2, localURL+"/"+taskId)
		h.checker(html2, localURL+"/"+taskId, taskId)
		visited = true
	})

	h.collector.Visit(h.ts.URL + localURL)

	return nil
}

func (h *Holberton) checker(html, url, taskId string) {
	c2 := colly.NewCollector()

	selector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div > div > div.modal-body > div.result", taskId)
	h.page.Click(selector)

	c2.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
	})

	c2.OnHTML(selector, func(element *colly.HTMLElement) {
		h4 := element.DOM.Find("h4")
		fmt.Println(h4.Text())

		selector := "div.check-inline"
		divs := element.DOM.Find(selector)
		divs.Each(func(_ int, check *goquery.Selection) {
			title, _ := check.Attr("title")
			fmt.Println(title)
		})
	})

	c2.Visit(h.ts.URL + url)
}
