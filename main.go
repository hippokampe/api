package main

import (
	"os"

	"github.com/hippokampe/api/app/api"
	"github.com/hippokampe/api/holberton"
	"github.com/hippokampe/api/logger"
)

func main() {
	hbtn, err := holberton.NewSession(holberton.FIREFOX)
	if err != nil {
		logger.Log2(err, "could not create the session")
		hbtn.CloseSession()
		os.Exit(1)
	}

	api.New(":5000", hbtn)

	hbtn.CloseSession()
}
