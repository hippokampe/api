package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/app/models"
	"github.com/hippokampe/api/holberton"
)

func getProjects(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Query("current") == "true" {
			projects, _ := h.GetCurrentProjects()
			ctx.JSON(http.StatusOK, gin.H{
				"current_projects": projects,
			})
			return
		}

		projects, _ := h.GetProjects()

		ctx.JSON(http.StatusOK, projects)
		if err := searcher.IndexProjects(projects.AllProjects); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
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

func searchProject(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitQuery := ctx.DefaultQuery("limit", "1")
		limit, err := strconv.Atoi(limitQuery)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		projectTitle, existsQuery := ctx.GetQuery("title")
		if !existsQuery || len(projectTitle) == 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "project title needs to be set",
			})
			return
		}

		if limit <= 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "limit must be greater than 0",
			})
			return
		}

		if limit == 1 {
			id, err := searcher.GetProjectID(projectTitle)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"project_id": id,
			})
		} else {
			results, err := searcher.GetProjects(projectTitle, limit)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, results)
		}

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

func logout(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := h.Logout(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Bye bye. See you soon"})
	}
}

func status(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "online",
		})
	}
}
