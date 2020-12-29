package holberton

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/utils"
)

func (hbtn *Holberton) getTask(taskElement *colly.HTMLElement) models.TaskBasic {
	var task models.TaskBasic

	// Get the position respect the others tasks
	taskPosition := utils.CleanString(taskElement.Attr("data-position"))

	// Search the ID
	idContainer := utils.CleanString(taskElement.Attr("data-role"))
	taskID := strings.TrimPrefix(idContainer, "task")

	// Task name
	taskTitle, taskType := parseTitleTask(taskElement.DOM)

	task.Position = taskPosition
	task.ID = taskID
	task.Title = taskTitle
	task.Type = taskType

	return task
}

func parseTitleTask(taskElementDOM *goquery.Selection) (title, class string) {
	h4 := taskElementDOM.Find("h4.task")
	span := taskElementDOM.Find("h4 > span")

	title = strings.Replace(h4.Text(), span.Text(), "", 1)
	title = utils.CleanString(title)
	class = utils.CleanString(span.Text())

	tmp := strings.SplitN(title, ".", 2)
	if len(tmp) == 2 {
		title = utils.CleanString(tmp[1])
	}

	if class[0] == '#' { //Ex, #advanced to advanced
		class = class[1:]
	}

	return title, class
}
