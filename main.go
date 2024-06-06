package main

import (
	"log"

	goauth "github.com/hanzalahimran7/go-auth/app"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app := goauth.Initialise()
	app.Run()
}
