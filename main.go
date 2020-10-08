package main

import (
	"github.com/hippokampe/configuration"
	"os"

	"github.com/hippokampe/api/app/api"
	"github.com/hippokampe/api/holberton"
	"github.com/hippokampe/api/logger"
)

func main() {
	config := configuration.New("/etc/hbtn/general.json")
	if err := config.ReadGeneralConfig(); err != nil {
		logger.Log2(err, "config generator")
		os.Exit(1)
	}

	hbtn, err := holberton.NewSession(holberton.FIREFOX, config)
	if err != nil {
		logger.Log2(err, "could not create the session")

		if hbtn != nil {
			hbtn.CloseSession()
		}
		os.Exit(1)
	}

	api.New(config.GetPort(), hbtn)

	hbtn.CloseSession()
}
