package holberton

import (
	"github.com/hippokampe/api/components/search"
	"github.com/hippokampe/api/models"
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
	ctx.Searcher = search.New("data/blave_resources", user.Email)
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

	if err := ctx.Searcher.IndexProjects(projects); err != nil {
		return models.Projects{}, errors.Wrap(err, "holberton: indexing")
	}

	return projects, nil
}

func (hbtn *Holberton) GetProject(email, id string) (models.Project, error) {
	scope := "holberton"
	ctx, err := hbtn.getSession(email)
	if err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	project, err := hbtn.getProject(*ctx.BrowserContext, id)
	if err != nil {
		return models.Project{}, errors.Wrap(err, scope)
	}

	return project, nil
}

func (hbtn *Holberton) SearchByTitle(email, title string, limit int) (interface{}, error) {
	scope := "holberton"
	ctx, err := hbtn.getSession(email)
	if err != nil {
		return "", errors.Wrap(err, scope)
	}

	if limit <= 0 {
		return "", ErrLimitNotValid
	}

	if limit == 1 {
		return ctx.Searcher.GetProjectID(title)
	}

	return ctx.Searcher.GetProjects(title, limit)
}
