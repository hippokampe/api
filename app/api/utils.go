package api

import (
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	logClient := logrus.New()

	//Disable logrus output
	/*src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err!= nil{
		fmt.Println("err", err)
	}*/

	logClient.Out = os.Stdout
	logClient.SetLevel(logrus.DebugLevel)
	apiLogPath := "api.log"
	logWriter, _ := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(apiLogPath),       // Generate a soft link to point to the latest log file
		rotatelogs.WithMaxAge(7*24*time.Hour),     // Maximum file save time
		rotatelogs.WithRotationTime(24*time.Hour), // Log cutting interval
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logClient.AddHook(lfHook)

	return func(c *gin.Context) {
		// Starting time
		start := time.Now()

		// Process request
		c.Next()

		// End Time
		end := time.Now()

		// Execution time
		latency := end.Sub(start) // latency = end - start

		path := c.Request.URL

		email, _ := c.Get("holberton_email")                 // Get the email with login endpoint
		if emailTmp, err := getEmailFromJWT(c); err == nil { // Get the email with the others endpoint
			email = emailTmp
		}

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		end = end.UTC()

		logClient.Infof("%s | %3d | %13v | %15s | %s  %s | +%v",
			email,
			statusCode,
			latency,
			clientIP,
			method, path,
			end,
		)
	}
}

func getEmailFromJWT(ctx *gin.Context) (string, error) {
	claims := jwt.ExtractClaims(ctx)
	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("type it's not compatible")
	}

	return email, nil
}
