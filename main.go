package main

import (
	"github.com/hippokampe/api/app/api"
	"log"

	"github.com/hippokampe/api/holberton"
)

func main() {
	hbtn, err := holberton.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := api.New(hbtn); err != nil {
		hbtn.Close()
		log.Fatal(err)
	}

	hbtn.Close()
}
