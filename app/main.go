package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	app := fiber.New()

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
