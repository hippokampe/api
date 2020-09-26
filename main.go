package main

import (
	"holberton/api/app/api"
	"holberton/api/holberton"
	"holberton/api/logger"
	"os"
)

/*func main() {
	var err error
	if err = holberton.NewSession(holberton.FIREFOX); err != nil {
		logger.Log2(err, "could not create the session")
		holberton.CloseSession()
		os.Exit(1)
	}

	holberton.StartPage()
	holberton.Login("1532@holbertonschool.com", "3006918Plata.")
	//holberton.GetProjects()
	//holberton.GetProject("314")
	holberton.CheckTask("304", "1776")
	holberton.CloseSession()

		1. Create daemon
		2. Start daemons as yuser

}
*/
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
