package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/holberton"
)

func getProjects(hbtn *holberton.Holberton) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		email, ok := claims["email"].(string)
		if !ok {
			ctx.AbortWithStatusJSON(500, gin.H{
				"message": "type email it's not the expected",
			})
			return
		}

		projects, err := hbtn.GetProjects(email)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})

			return
		}

		ctx.JSON(200, gin.H{
			"projects": projects,
			"total":    len(projects),
		})

	}
}
