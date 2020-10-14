package holberton

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/logger"
	"sort"
	"strings"
)

func (h *Holberton) projects() (models.Projects, error) {
	var err error
	var projects models.Projects

	_, err = h.page.Goto("https://intranet.hbtn.io/dashboards/my_current_projects")
	if err != nil {
		logger.Log2(err, "could not goto")
		return models.Projects{}, err
	}

	html, _ := h.page.Content()
	url := h.setHtml(html, "/projects")

	h.collector.OnHTML("body > main > article", func(element *colly.HTMLElement) {
		projects.CurrentProjects = searchCurrentProjects(element)

		selector := "div.panel.panel-default"
		element.ForEach(selector, func(_ int, e *colly.HTMLElement) {
			h4 := e.DOM.Find("h4.panel-title")
			title := cleanString(h4.Text())

			projectsList := e.DOM.Find("li.list-group-item")
			projectsList.Each(func(_ int, project *goquery.Selection) {
				name := project.Find("a")       // Project title
				code := project.Find("code")    // Project code or id
				score := project.Find("strong") // Score of the project

				projects.AllProjects = append(projects.AllProjects, models.Project{
					Title:    name.Text(),
					ID:       code.Text(),
					Score:    score.Text(),
					Category: title,
				})
			})
		})

	})

	_ = h.collector.Visit(url)

	return projects, nil
}

func (h *Holberton) currentProjects() (models.CurrentProjects, error) {
	var err error
	var currentProjects models.CurrentProjects

	_, err = h.page.Goto("https://intranet.hbtn.io/dashboards/my_current_projects")
	if err != nil {
		logger.Log2(err, "could not goto")
		return models.CurrentProjects{}, err
	}

	html, _ := h.page.Content()
	url := h.setHtml(html, "/projects")

	h.collector.OnHTML("body > main > article", func(element *colly.HTMLElement) {
		currentProjects = searchCurrentProjects(element)
	})

	_ = h.collector.Visit(url)

	return currentProjects, nil
}

func searchCurrentProjects(html *colly.HTMLElement) models.CurrentProjects {
	var currentProjects models.CurrentProjects
	var secondDeadline []models.Project
	var firstDeadline []models.Project

	selector := "ul.list-group:nth-child(2)"
	selection := html.DOM.Find(selector)
	list := selection.Children()
	list.Each(func(_ int, selection *goquery.Selection) {
		title := selection.Find("a")
		category := selection.Find("em")
		code := selection.Find("code")

		scoreSelector := selection.Find("div.project_progress_percentage.alert.alert-info")
		score := scoreSelector.Text()
		score = cleanString(score)
		score = strings.TrimSuffix(score, " done")

		// Deadline summary
		deadlineGeneral := selection.Find("span.bpi-status").Text()
		information := strings.Split(deadlineGeneral, ", ")
		period, finished, remaining := getDeadline(information[1])
		startedTmp := strings.Split(deadlineGeneral, ",")[0]
		started := strings.TrimPrefix(startedTmp, "started on ")

		project := models.Project{
			ID:       code.Text(),
			Title:    title.Text(),
			Score:    score,
			Category: category.Text(),
			DeadlineInformation: models.DeadlineSummary{
				Period:        period,
				Started:       started,
				Finished:      finished,
				RemainingDate: remaining,
			},
		}

		if project.DeadlineInformation.Period == 1 {
			firstDeadline = append(firstDeadline, project)
		} else {
			secondDeadline = append(secondDeadline, project)
		}
	})

	sort.Slice(firstDeadline, func(i, j int) bool {
		return firstDeadline[i].DeadlineInformation.RemainingDate < firstDeadline[j].DeadlineInformation.RemainingDate
	})

	sort.Slice(secondDeadline, func(i, j int) bool {
		return secondDeadline[i].DeadlineInformation.RemainingDate < secondDeadline[j].DeadlineInformation.RemainingDate
	})

	currentProjects.FirstDeadline = firstDeadline
	currentProjects.SecondDeadline = secondDeadline
	currentProjects.Total = len(firstDeadline) + len(secondDeadline)

	return currentProjects
}
