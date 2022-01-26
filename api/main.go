package main

import (
	"fmt"
	"log"
	"os"
	"shorten-url-with-redis/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setUpRoute(app *fiber.App) {
	app.Get("/:url", route.ResovleUrl)
	app.Post("/api/v1", route.ShortenUrl)
}

func main() {
	err := godoenv.Load()

	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()
	app.Use(logger.New())
	setUpRoute(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
