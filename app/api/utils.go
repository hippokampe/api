package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func getEmailFromJWT(ctx *gin.Context) (string, error) {
	claims := jwt.ExtractClaims(ctx)
	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("type it's not compatible")
	}

	return email, nil
}
