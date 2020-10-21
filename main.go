package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	general "github.com/hippokampe/configuration/v2"
	"github.com/hippokampe/configuration/v2/configuration"
	"github.com/hippokampe/configuration/v2/credentials"

	"github.com/hippokampe/api/app/api"
	"github.com/hippokampe/api/holberton"
	"github.com/hippokampe/api/logger"
)

var (
	cred   *credentials.Credentials
	config *configuration.InternalSettings
)

func browserSelector(browserPath string) string {
	if strings.Contains(browserPath, "firefox") {
		return holberton.FIREFOX
	}

	if strings.Contains(browserPath, "chromium") {
		return holberton.CHROMIUN
	}

	if strings.Contains(browserPath, "webkit") {
		return holberton.WEBKIT
	}

	return ""
}

func init() {
	fmt.Println("init")


	gin.SetMode(gin.ReleaseMode)

	generalConfig := general.New("/etc/hippokampe/general.json")
	if err := generalConfig.ReadGeneralConfig(); err != nil {
		log.Fatal("cannot read the /etc/hippokampe/general.json. Check the documentation")
	}

	if err := os.Setenv("HIPPOKAMPE_CONFIGURATION", generalConfig.CustomSettingsFilename); err != nil {
		log.Fatal("cannot set the $HIPPOKAMPE_CONFIGURATION. Check the documentation")
	}

	if err := os.Setenv("HIPPOKAMPE_CREDENTIALS", generalConfig.CredentialsFilename); err != nil {
		log.Fatal("cannot set the $HIPPOKAMPE_CREDENTIALS. Check the documentation")
	}

	config = configuration.New()
	cred = credentials.New()

	if err := config.BindCredentials(cred); err != nil {
		log.Fatal(err)
	}

	config.SetFilename(os.Getenv("HIPPOKAMPE_CONFIGURATION"))
	cred.SetFilename(os.Getenv("HIPPOKAMPE_CREDENTIALS"))

	if err := config.ReadFromFile(); err != nil {
		log.Fatal(err)
	}

	_ = cred.ReadFromFile()
}

func main() {
	fmt.Println("main")


	browserPath, _ := config.GetPathBrowser()
	hbtn, err := holberton.NewSession(browserSelector(browserPath), config)
	if err != nil {
		logger.Log2(err, "could not create the session")

		if hbtn != nil {
			hbtn.CloseSession()
		}
		os.Exit(1)
	}

	if err := api.New(hbtn, config); err != nil {
		if hbtn != nil {
			hbtn.CloseSession()
		}

		log.Fatal(err)
	}

	hbtn.CloseSession()
}
