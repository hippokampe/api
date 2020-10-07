package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/holberton"
)

func New(port string, holberton *holberton.Holberton) {
	router := gin.Default()

	router.Use(StatusApp(holberton))

	router.GET("/status", status(holberton))
	router.POST("/login", login(holberton))

	authorized := router.Group("/")
	authorized.Use(Authorized(holberton))
	{
		authorized.GET("/projects", getProjects(holberton))
		authorized.GET("/projects/:id", getProject(holberton))
		authorized.GET("/projects/:id/checker/:task", checkTask(holberton))
	}

	router.Run(port)
}
