package log

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func init() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("log.txt open fail.")
	}

	Info = log.New(file, "info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(file, "warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "error: ", log.Ldate|log.Ltime|log.Lshortfile)
}
