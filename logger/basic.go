package logger

import (
	"errors"
	"log"
)

func Log(err error) {
	log.Fatal(err)
}

func Log2(err error, message string) {
	log.Fatalf("%s: %s\n", message, err.Error())
}

func New(message string) error {
	return errors.New(message)
}
