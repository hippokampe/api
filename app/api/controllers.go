package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/holberton"
	"github.com/hippokampe/api/models"
	"github.com/pkg/errors"
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

func checker(hbtn *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email, err := getEmailFromJWT(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "type email it's not the expected",
			})
			return
		}

		projectID := ctx.Param("id")
		taskID := ctx.Param("task")

		if len(projectID) != 3 || len(taskID) != 4 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "the id submitted for the project or the task it's not valid",
			})
			return
		}

		taskChecker, err := hbtn.Checker(email, projectID, taskID, false)
		if err != nil {
			var errorMessage string
			var errorCode int
			switch errors.Cause(err) {
			case holberton.ErrType:
				errorMessage = holberton.ErrType.Error()
				errorCode = http.StatusBadRequest
			case holberton.ErrTaskNotFound:
				errorMessage = holberton.ErrTaskNotFound.Error()
				errorCode = http.StatusNotFound
			default:
				errorMessage = err.Error()
				errorCode = http.StatusInternalServerError
			}

			ctx.JSON(errorCode, gin.H{
				"message": errorMessage,
			})
			return
		}

		ctx.JSON(http.StatusOK, taskChecker)
	}
}

func checker2(hbtn *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email, err := getEmailFromJWT(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "type email it's not the expected",
			})
			return
		}

		projectID := ctx.Param("id")
		taskPosition := ctx.Query("position")

		if len(projectID) != 3 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "the id submitted for the project it's not valid",
			})
			return
		}

		taskChecker, err := hbtn.Checker(email, projectID, taskPosition, true)
		if err != nil {
			var errorMessage string
			var errorCode int
			switch errors.Cause(err) {
			case holberton.ErrType:
				errorMessage = holberton.ErrType.Error()
				errorCode = http.StatusBadRequest
			case holberton.ErrTaskNotFound:
				errorMessage = holberton.ErrTaskNotFound.Error()
				errorCode = http.StatusNotFound
			default:
				errorMessage = err.Error()
				errorCode = http.StatusInternalServerError
			}

			ctx.JSON(errorCode, gin.H{
				"message": errorMessage,
			})
			return
		}

		ctx.JSON(http.StatusOK, taskChecker)
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

func searchProject(hbtn *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email, err := getEmailFromJWT(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "type email it's not the expected",
			})
			return
		}

		limitQuery := ctx.DefaultQuery("limit", "1")
		limit, err := strconv.Atoi(limitQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		projectTitle, existsQuery := ctx.GetQuery("title")
		if !existsQuery || len(projectTitle) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "project title needs to be set",
			})
			return
		}

		result, err := hbtn.SearchByTitle(email, projectTitle, limit)
		if err != nil {
			switch errors.Cause(err) {
			case holberton.ErrLimitNotValid:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}

			return
		}

		resultSearch, ok := result.(models.ProjectsResultSearch)
		if !ok {
			projectId := result.(string)
			ctx.JSON(http.StatusOK, gin.H{
				"project_id": projectId,
			})
			return
		}

		ctx.JSON(http.StatusOK, resultSearch)
	}
}
