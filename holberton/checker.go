package holberton

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/hippokampe/api/models"
	"github.com/hippokampe/api/utils"
	"github.com/mxschmitt/playwright-go"
	"github.com/pkg/errors"
)

func (hbtn *Holberton) checkByTaskID(browserCtx playwright.BrowserContext, projectID, taskID string) (models.TaskChecker, error) {
	scope := "checker by task id"
	page, err := browserCtx.NewPage()
	if err != nil {
		return models.TaskChecker{}, errors.Wrap(err, scope)
	}

	defer page.Close()

	projectUrl := fmt.Sprintf("https://intranet.hbtn.io/projects/%s", projectID)
	if _, err := page.Goto(projectUrl); err != nil {
		return models.TaskChecker{}, errors.Wrap(err, scope)
	}

	baseTaskSelector := fmt.Sprintf("#task-%s", taskID)
	checkBtnSelector := fmt.Sprintf("%s > div.student_correction_requests > button.task_correction_modal.btn.btn-default", baseTaskSelector)
	newTestBtnSelector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div > div"+
		" > div.modal-body > div.actions > center > input", taskID)
	resultSelector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div"+
		" > div > div.modal-body > div.result", taskID)

	if page.Click(checkBtnSelector) != nil {
		return models.TaskChecker{}, errors.Wrap(ErrTaskNotFound, scope)
	}

	if page.Click(newTestBtnSelector) != nil {
		return models.TaskChecker{}, errors.Wrap(ErrTaskNotFound, scope)
	}

	if _, err := page.WaitForSelector(resultSelector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Int(1000 * 60),
	}); err != nil {
		return models.TaskChecker{}, errors.Wrap(err, scope)
	}

	html, err := page.Content()
	if err != nil {
		return models.TaskChecker{}, errors.Wrap(err, scope)
	}

	url := hbtn.setHtml(html)

	resultDisplay, err := hbtn.extractChecks(url, taskID)
	if err != nil {
		return models.TaskChecker{}, errors.Wrap(err, scope)
	}

	var task models.TaskChecker

	url = hbtn.setHtml(html)
	hbtn.collector.OnHTML(baseTaskSelector, func(element *colly.HTMLElement) {
		tmp := hbtn.getTask(element.DOM.Parent())
		task.ID = tmp.ID
		task.Title = tmp.Title
		task.Position = tmp.Position
		task.Type = tmp.Type
	})

	if err := hbtn.collector.Visit(url); err != nil {
		return models.TaskChecker{}, nil
	}

	task.ResultDisplay = resultDisplay

	return task, nil
}

func (hbtn *Holberton) checkByIndex(browserCtx playwright.BrowserContext, projectID string, taskIndex int) (models.TaskChecker, error) {
	scope := "checker by index"
	page, err := browserCtx.NewPage()
	if err != nil {
		return models.TaskChecker{}, errors.Wrap(err, scope)
	}

	defer page.Close()

	visited := false
	found := false

	projectUrl := fmt.Sprintf("https://intranet.hbtn.io/projects/%s", projectID)
	if _, err := page.Goto(projectUrl); err != nil {
		return models.TaskChecker{}, errors.Wrap(err, scope)
	}

	html, err := page.Content()
	if err != nil {
		return models.TaskChecker{}, errors.Wrap(err, scope)
	}

	var taskDocument *goquery.Document

	url := hbtn.setHtml(html)
	hbtn.collector.OnHTML("main", func(sectionElement *colly.HTMLElement) {
		if visited {
			return
		}

		section := sectionElement.DOM.Find("article > section.formatted-content").Last()
		tasksSelection := section.Children()

		if taskIndex > tasksSelection.Length() {
			return
		}

		taskSelected := tasksSelection.Get(taskIndex - 1)
		taskDocument = goquery.NewDocumentFromNode(taskSelected)

		visited = true
		found = true
	})

	if err := hbtn.collector.Visit(url); err != nil {
		return models.TaskChecker{}, nil
	}

	if !found || taskDocument == nil {
		return models.TaskChecker{}, errors.Wrap(ErrTaskNotFound, scope)
	}

	val, exists := taskDocument.Attr("data-role")
	if !exists {
		return models.TaskChecker{}, errors.Wrap(ErrTaskNotFound, scope)
	}

	taskID := strings.TrimPrefix(val, "task")
	return hbtn.checkByTaskID(browserCtx, projectID, taskID)
}

func (hbtn *Holberton) extractChecks(url, taskID string) (models.ResultDisplay, error) {
	scope := "extract checks"

	checksSelector := fmt.Sprintf("#task-test-correction-%s-correction-modal > div > div > "+
		"div.modal-body > div.result", taskID)

	var resultDisplay models.ResultDisplay
	hbtn.collector.OnHTML(checksSelector, func(resultElement *colly.HTMLElement) {
		done := true
		resultElement.ForEach("div.check-inline", func(_ int, checkElement *colly.HTMLElement) {
			titleAttr := checkElement.Attr("title")
			checkLabel, passed := parseCheckTitle(titleAttr)
			if !passed {
				done = false
			}

			title := checkElement.Text
			title = utils.CleanString(title)

			resultDisplay.Checks = append(resultDisplay.Checks, models.Check{
				Title:      title,
				CheckLabel: checkLabel,
				Passed:     passed,
			})
		})

		resultDisplay.Done = done
	})

	if err := hbtn.collector.Visit(url); err != nil {
		return models.ResultDisplay{}, errors.Wrap(err, scope)
	}

	return resultDisplay, nil
}

func parseCheckTitle(titleAttr string) (checkLabel string, status bool) {
	titleAttr = utils.CleanString(titleAttr)
	titleAttr = strings.ToLower(titleAttr)

	tmp := strings.Split(titleAttr, " - ")
	tmp2 := strings.Split(tmp[0], " ")
	if len(tmp2) >= 2 {
		checkLabel = tmp2[len(tmp2)-1]
	} else {
		checkLabel = tmp[0]
	}

	status = false
	if tmp[1] == "success" {
		status = true
	}

	return checkLabel, status
}
