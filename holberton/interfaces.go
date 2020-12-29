package holberton

import (
	"github.com/hippokampe/api/app/models"
	"github.com/pkg/errors"
)

func (hbtn *Holberton) Login(credentials models.Login) (models.User, error) {
	scope := "holberton"
	ctx, err := hbtn.getSession(credentials.Email)
	if err == nil {
		return *ctx.User, nil
	}

	ctx, err = hbtn.newSession()
	if err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	user, err := hbtn.login(*ctx.BrowserContext, credentials)
	if err != nil {
		return models.User{}, errors.Wrap(err, scope)
	}

	ctx.User = &user
	hbtn.addSession(user.Email, ctx)
	return user, nil
}

func (hbtn *Holberton) GetProjects(email string) (models.Projects, error) {
	scope := "holberton"
	ctx, err := hbtn.getSession(email)
	if err != nil {
		return models.Projects{}, errors.Wrap(err, scope)
	}

	projects, err := hbtn.getProjects(*ctx.BrowserContext)
	if err != nil {
		return models.Projects{}, errors.Wrap(err, scope)
	}

	return projects, nil
}

func (hbtn *Holberton) GetProject(email, id string) (models.Project, error) {
	scope := "holberton"
	ctx, err := hbtn.getSession(email)
	if err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	projects, err := hbtn.getProject(*ctx.BrowserContext, id)
	if err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	return projects, nil
}
