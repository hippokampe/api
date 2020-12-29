package api

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hippokampe/api/holberton"
	"github.com/hippokampe/api/models"
)

func JWTMiddleware(hbtn *holberton.Holberton) (*jwt.GinJWTMiddleware, error) {
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:              []byte("hippokampe"),
		PrivKeyFile:      "keys/jwtRS256.key",
		PubKeyFile:       "keys/jwtRS256.key.pub",
		SigningAlgorithm: "RS256",
		Realm:            "zone",
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		IdentityKey:      "email",
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			return authenticator(hbtn, ctx)
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return payloadGenerator(data)
		},
		IdentityHandler: func(ctx *gin.Context) interface{} {
			return identityHandler(ctx)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.User); ok && hbtn.IsLogged(v.Email) {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	return middleware, err
}

func authenticator(hbtn *holberton.Holberton, ctx *gin.Context) (interface{}, error) {
	var credentials models.Login
	if err := ctx.ShouldBind(&credentials); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	user, err := hbtn.Login(credentials)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return user, nil
}

func payloadGenerator(data interface{}) jwt.MapClaims {
	if v, ok := data.(models.User); ok {
		return jwt.MapClaims{
			"email": v.Email,
		}
	}

	return jwt.MapClaims{}
}

func identityHandler(ctx *gin.Context) interface{} {
	claims := jwt.ExtractClaims(ctx)
	return &models.User{
		Email: claims["email"].(string),
	}
}
