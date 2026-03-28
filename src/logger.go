package main


import (
	"log"
	"os"
)



func InitLogger(){
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(logfile)
}