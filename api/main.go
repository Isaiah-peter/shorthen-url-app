package main

import (
	"shorten-url-with-redis/route"

	"github.com/gofiber/fiber/v2"
)

func setUpRoute(app *fiber.App) {
	app.Get("/:url", route.ResovleUrl)
	app.Post("/api/v1", route.ShortenUrl)
}

func main() {
	app := fiber.New()

}
