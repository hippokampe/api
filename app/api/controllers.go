package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/holberton"
)

func getProjects(hbtn *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email, err := getEmailFromJWT(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "type email it's not the expected",
			})
			return
		}

		projects, err := hbtn.GetProjects(email)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"projects": projects,
			"total":    len(projects),
		})

	}
}

func getProject(hbtn *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email, err := getEmailFromJWT(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "type email it's not the expected",
			})
			return
		}

		projectId := ctx.Param("id")
		if len(projectId) != 3 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "the id param it's not valid",
			})
			return
		}

		project, err := hbtn.GetProject(email, projectId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		if project.Exists() {
			ctx.JSON(http.StatusOK, project)
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "project not found",
		})
	}
}
