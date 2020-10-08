package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/holberton"
	"github.com/hippokampe/configuration"
)

var config *configuration.Configuration

func New(holberton *holberton.Holberton, configParam *configuration.Configuration) error {
	router := gin.Default()

	config = configParam

	router.GET("/status", status(holberton))
	router.POST("/login", login(holberton))

	authorized := router.Group("/")
	authorized.Use(Authorized(holberton))
	{
		authorized.GET("/projects", getProjects(holberton))
		authorized.GET("/projects/:id", getProject(holberton))
		authorized.GET("/projects/:id/checker/:task", checkTask(holberton))
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
	
	if config.InternalStatus.Logged {
		email := config.InternalStatus.Credentials.Email
		password := config.InternalStatus.Credentials.Password
		_, _ = holberton.Login(email, password)
	}

	return nil
}