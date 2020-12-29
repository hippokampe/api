package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/holberton"
	"github.com/pkg/errors"
)

func New(hbtn *holberton.Holberton) error {
	scope := "api: new"
	router := gin.Default()

	authMiddleware, err := JWTMiddleware(hbtn)
	if err != nil {
		return errors.Wrap(err, scope)
	}

	if err := authMiddleware.MiddlewareInit(); err != nil {
		return errors.Wrap(err, scope)
	}

	router.POST("/auth/login", authMiddleware.LoginHandler)
	contextHandler := router.Group("/")
	contextHandler.Use(authMiddleware.MiddlewareFunc())
	{
		contextHandler.GET("/auth/refresh_token", authMiddleware.RefreshHandler)

		contextHandler.GET("/projects", getProjects(hbtn))
		contextHandler.GET("/projects/:id", getProject(hbtn))

		contextHandler.GET("/search", searchProject(hbtn))
	}

	return router.Run(":8080")
}
