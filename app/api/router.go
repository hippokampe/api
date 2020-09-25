package api

import (
	"fmt"
	"holberton/api/app/models"
	"holberton/api/holberton"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(port string, holberton holberton.Holberton) {
	router := gin.Default()

	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Alive",
		})
	})

	router.POST("/login", func(c *gin.Context) {
		var user models.User

		fmt.Println("logging")
		fmt.Println(c.Param("email"))
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(user.Email)

		holberton.StartPage()
		newUser, _ := holberton.Login(user.Email, user.Password)

		c.JSON(http.StatusOK, gin.H{"username": newUser.Username})
	})

	router.GET("/projects", func(c *gin.Context) {
		holberton.StartPage()
		projects, _ := holberton.GetProjects()
		c.JSON(http.StatusOK, projects)
	})

	router.GET("/projects/:id", func(c *gin.Context) {
		id := c.Param("id")

		project, _ := holberton.GetProject(id)

		c.JSON(http.StatusOK, project)
	})

	router.GET("/projects/:id/checker/:task", func(c *gin.Context) {
		id := c.Param("id")
		taskID := c.Param("task")

		holberton.CheckTask(id, taskID)
	})

	router.Run(port)
}
