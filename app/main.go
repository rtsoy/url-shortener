package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rtsoy/url-shortener/routes"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	app := fiber.New()

	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
