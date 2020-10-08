package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/holberton"
)

func getProjects(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		projects, _ := h.GetProjects()
		ctx.JSON(http.StatusOK, projects)
	}
}

func getProject(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		project, _ := h.GetProject(id)
		if project == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		}

		ctx.JSON(http.StatusOK, project)
	}
}

func checkTask(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		taskID := ctx.Param("task")

		task, err := h.CheckTask(id, taskID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, task)
	}
}

func login(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if h.InternalStatus.Logged {
			ctx.JSON(http.StatusOK, gin.H{"username": h.InternalStatus.Username})
			return
		}

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser, err := h.Login(user.Email, user.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if newUser == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "user email or user password is incorrect",
			})
			return
		}

		h.InternalStatus.Username = newUser.Username

		ctx.JSON(http.StatusOK, gin.H{"username": newUser.Username})
	}
}

func status(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "online",
		})
	}
}
