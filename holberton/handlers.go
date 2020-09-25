package holberton

import (
	"holberton/api/app/models"
	"holberton/api/logger"
)

func (h *Holberton) StartPage() error {
	var err error

	if h.page != nil {
		return nil
	}

	h.page, err = h.browser.NewPage()
	if err != nil {
		logger.Log2(err, "could not create page")
		return err
	}

	h.newServer()
	return nil
}

func (h *Holberton) Login(email, password string) (*models.User, error) {
	return h.login(email, password)
}

func (h *Holberton) GetProjects() (models.Projects, error) {
	return h.projects()
}

func (h *Holberton) GetProject(id string) (*models.Project, error) {
	return h.project(id)
}

func (h *Holberton) CheckTask(id, taskId string) error {
	return h.checkTask(id, taskId)
}
