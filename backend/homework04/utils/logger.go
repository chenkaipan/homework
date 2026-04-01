package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogger() {

	InfoLogger = log.New(
		os.Stdout,
		"INFO\t",
		log.LstdFlags|log.Lshortfile,
	)

	ErrorLogger = log.New(
		os.Stderr,
		"ERROR\t",
		log.LstdFlags|log.Lshortfile,
	)
}
