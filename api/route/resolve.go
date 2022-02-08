package route

import (
	"context"
	"shorten-url-with-redis/database"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResovleUrl(c *fiber.Ctx) error {
	var ctx = context.Background()
	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "url does not exist"})
	} else if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "error has occur"})
	}

	rimr := database.CreateClient(1)
	defer rimr.Close()
	_ = rimr.Incr(ctx, "counter")
	return c.Redirect(value, 301)
}
