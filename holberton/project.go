package holberton

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/logger"
	"log"
	"os"
	"path/filepath"
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
	url := h.setHtml(html, localURL)

	h.collector.OnHTML("article", func(article *colly.HTMLElement) {
		if visited {
			return
		}

		categoryP := article.DOM.Find("body > main > article > p.sm-gap")
		categoryTitle := cleanString(categoryP.Text())
		project.Category = categoryTitle

		projectTitle := article.DOM.Find("h1.gap")
		project.Title = projectTitle.Text()

		_, err := h.generateReadme(project, article.DOM.Find("article"))
		if err != nil {
			log.Println(err)
		}

		taskPath, err := h.createDirTasks(project.Title)
		if err != nil {
			log.Fatal(err)
			return
		}

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

			h4, span := searchTitleTask(task)
			title, class := parseTitleTask(h4, span)
			status := searchTaskDone(task)

			taskDescription, err := h.generateTask(taskPath, title, task)
			if err != nil {
				log.Println(err)
				return
			}

			project.Tasks = append(project.Tasks, models.Task{
				ID:              taskID,
				Title:           title,
				Type:            class,
				Done:            status,
				FileDescription: taskDescription,
			})
		})

		visited = true
	})

	_ = h.collector.Visit(url)

	if project.Title == "" {
		return nil, nil
	}

	return project, nil
}

func (h *Holberton) generateReadme(project *models.Project, selection *goquery.Selection) (string, error) {
	filename := filepath.Join("/home/davixcky", project.Title) + ".md"
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	selection.Children().Find("a").Each(func(_ int, el *goquery.Selection) {
		link, _ := el.Attr("href")
		link = "https://intranet.hbtn.io" + link
		el.SetAttr("href", link)
	})

	converter := md.NewConverter("", true, nil)
	markdown := converter.Convert(selection)

	// Writing title of the file (project title)
	title := fmt.Sprintf("# %s\n", project.Title)
	_, err = file.WriteString(title)
	if err != nil {
		return "", err
	}

	// Writing content of the project
	_, err = file.WriteString(markdown)
	if err != nil {
		return "", err
	}

	if err := file.Sync(); err != nil {
		return "", nil
	}

	return filename, nil
}
