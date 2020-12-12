package api

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/app/components/search"
	"github.com/hippokampe/api/holberton"
	"github.com/hippokampe/configuration/v2/configuration"
)

var config *configuration.InternalSettings
var searcher *search.Search

func New(holberton *holberton.Holberton, configParam *configuration.InternalSettings) error {
	router := gin.Default()

	config = configParam
	dir, _ := filepath.Split(os.Getenv("HIPPOKAMPE_CONFIGURATION"))
	filename := filepath.Join(dir, "projects")
	searcher = search.New(filename)

	router.GET("/status", status(holberton))
	router.POST("/login", login(holberton))

	authorized := router.Group("/")
	authorized.Use(Authorized(holberton))
	{
		authorized.POST("/logout", logout(holberton))
		authorized.GET("/projects", getProjects(holberton))
		authorized.GET("/projects/:id", getProject(holberton))
		authorized.GET("/projects/:id/checker/:task", checkTask(holberton))
		authorized.GET("/search", searchProject(holberton))
	}

	if err := restoreSession(holberton); err != nil {
		return err
	}

	return router.Run(config.GetPort())
}

func restoreSession(holberton *holberton.Holberton) error {
	if err := holberton.StartPage(); err != nil {
		return err
	}

	status, err := config.IsLogged()
	if err != nil {
		return err
	}

	if status { // Restore the previous session
		cred, err := config.GetCredentials()
		if err != nil {
			return err
		}

		email, _ := cred.GetValue("email")
		password, _ := cred.GetValue("password")
		if _, err = holberton.Login(email, password); err != nil {
			return err
		}

		projects, err := holberton.GetProjects()
		if err != nil {
			return err
		}

		return searcher.IndexProjects(projects.AllProjects)
	}

	return nil
}
