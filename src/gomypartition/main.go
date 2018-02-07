package main

import (
	"gomypartition/Application"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	//loading configuration file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app:= new(Application.Application)
	app.Run()
}
