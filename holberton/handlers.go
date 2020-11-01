package holberton

import (
	"errors"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/logger"
	"strings"
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

	h.InternalStatus.VisitedURLS = make(map[string]bool)
	h.InternalStatus.Logged = false

	return nil
}

func (h *Holberton) Login(email, password string) (*models.User, error) {
	user, err := h.login(email, password)
	if err != nil {
		if strings.Contains(err.Error(), "waiting for selector \"#new_user > div.actions > input\"") {
			err = errors.New("bad credentials")
		}

		return nil, err
	}

	return user, nil
}

func (h *Holberton) Logout() error {
	return h.logout()
}

func (h *Holberton) GetProjects() (models.Projects, error) {
	return h.projects()
}

func (h *Holberton) GetProject(id string) (*models.Project, error) {
	return h.project(id)
}

func (h *Holberton) GetCurrentProjects() (models.CurrentProjects, error) {
	return h.currentProjects()
}

func (h *Holberton) CheckTask(id, taskId string) (*models.Task, error) {
	return h.checkTask(id, taskId)
}
