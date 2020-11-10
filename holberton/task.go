package holberton

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func (h *Holberton) generateTask(taskPath, taskFilename, titleTask string, selection *goquery.Selection) (string, error) {
	// Cleaning html
	selection.Find("div.task_progress_score_bar").Remove()
	selection.Find("div.student_task_controls").Remove()
	selection.Find("div.student_correction_requests").Remove()
	selection.Find("h4.task").Remove()

	filename := filepath.Join(taskPath, taskFilename) + ".md"
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
	basicPath := os.Getenv("HIPPOKAMPE")
	path := filepath.Join(basicPath, "projects", titleProject)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	group, err := user.LookupGroup("hippokampe")
	if err != nil {
		return "", err
	}

	uid, _ := strconv.Atoi(usr.Uid)
	gid, _ := strconv.Atoi(group.Gid)

	if err := os.Chown(path, uid, gid); err != nil {
		return "", err
	}

	fmt.Println(group.Name)

	return path, nil
}
