package holberton

import (
	"fmt"
	"os"
	"path/filepath"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

func (h *Holberton) generateTask(taskPath, titleTask string, selection *goquery.Selection) (string, error) {
	// Cleaning html
	selection.Find("div.task_progress_score_bar").Remove()
	selection.Find("div.student_task_controls").Remove()
	selection.Find("div.student_correction_requests").Remove()
	selection.Find("h4.task").Remove()

	filename := filepath.Join(taskPath, titleTask) + ".md"
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	defer file.Close()

	fixHolbertonLinks(selection)

	converter := md.NewConverter("", true, nil)
	markdown := converter.Convert(selection)

	// Writing title of the file (task title)
	title := fmt.Sprintf("# %s\n", titleTask)
	_, err = file.WriteString(title)
	if err != nil {
		return "", err
	}

	// Writing content of the task
	_, err = file.WriteString(markdown)
	if err != nil {
		return "", err
	}

	return filename, file.Sync()
}

func (h *Holberton) createDirTasks(titleProject string) (string, error) {
	basicPath := h.Configuration.InternalStatus.ConfigurationFile
	path := filepath.Join(basicPath, "tasks", titleProject)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}

	return path, nil
}
