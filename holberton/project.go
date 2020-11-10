package holberton

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/logger"
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

		taskPath, err := h.createDirTasks(project.ID)
		if err != nil {
			log.Fatal(err)
			return
		}

		_, err = h.generateReadme(project, article.DOM.Find("article"))
		if err != nil {
			log.Println(err)
		}

		scoreSelector := "body > main > article > div.gap.clean.well > ul > li:nth-child(2) >" +
			" ul > li:nth-child(3) > strong"
		projectScore := article.DOM.Find(scoreSelector)
		project.Score = projectScore.Text()

		taskIndex := 0 // Task index

		tasksContainerSelector := "body > main > article > section"
		tasksContainer := article.DOM.Find(tasksContainerSelector)
		tasksContainer.Children().Each(func(_ int, taskSelector *goquery.Selection) {
			if taskSelector.Is("div.quiz_question_item_container") {
				return
			}

			value, _ := taskSelector.Attr("data-role")
			taskID := value[4:]

			h4, span := searchTitleTask(taskSelector)
			title, class := parseTitleTask(h4, span)
			status := searchTaskDone(taskSelector)

			index := strconv.Itoa(taskIndex)
			taskDescription, err := h.generateTask(taskPath, index, title, taskSelector)
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

			taskIndex++
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
	basicPath := os.Getenv("HIPPOKAMPE")
	filename := filepath.Join(basicPath, "projects", project.ID, "basic_information.md")
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	fixHolbertonLinks(selection)

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

	if err := chownR(filename); err != nil {
		return "", err
	}

	return filename, nil
}
