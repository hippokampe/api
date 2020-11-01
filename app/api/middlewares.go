package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/holberton"
)

func Authorized(h *holberton.Holberton) gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.InternalStatus.Logged {
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "you must have logged in",
		})
	}
}

func StatusApp(h *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if h.InternalStatus.Started {
			return
		}

		if ctx.FullPath() == "/status" {
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "scrapper has not started. Check documentation",
		})
	}
}
