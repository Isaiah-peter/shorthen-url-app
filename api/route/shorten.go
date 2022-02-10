package route

import (
	"context"
	"os"
	"strconv"
	"time"

	"shorten-url-with-redis/database"
	"shorten-url-with-redis/helper"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_limit"`
	XRateLimitRest time.Duration `json:"rate_limit_rest"`
}

func ShortenUrl(c *fiber.Ctx) error {
	ctx := context.Background()
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse json"})
	}
	//implement rate limiting

	r2 := database.CreateClient(1)
	defer r2.Close()
	value, err := r2.Get(ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(ctx, c.IP(), os.Getenv("API_QOUTA"), 30*60*time.Second)
	} else {
		value, _ := r2.Get(database.Ctx, c.IP()).Result()
		valueInt, _ := strconv.Atoi(value)
		if valueInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "rate exceeded", "rate_limit_rest": limit / time.Nanosecond / time.Minute})
		}

	}

	//check if the input sent is a URL
	if !govalidator.isURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "this input is not an URL"})
	}

	//check for domain error
	if !helper.RemoveDomainingError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "you cannot use this url"})
	}
	//enforce  http SSL

	body.URL = helper.EnforceHTTP(body.URL)
	r2.Decr(database.Ctx, c.IP())
}
