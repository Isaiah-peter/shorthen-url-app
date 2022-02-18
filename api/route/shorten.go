package route

import (
	"context"
	"os"
	"strconv"
	"time"

	"shorten-url-with-redis/database"
	"shorten-url-with-redis/helper"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_rest"`
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
		valueInt, _ := strconv.Atoi(value)
		if valueInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "rate exceeded", "rate_limit_rest": limit / time.Nanosecond / time.Minute})
		}

	}

	//check if the input sent is a URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "this input is not an URL"})
	}

	//check for domain error
	if !helper.RemoveDomainingError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "you cannot use this url"})
	}
	//enforce  http SSL

	body.URL = helper.EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ := r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fail this url is in use"})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	errd := r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second)

	if errd == nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unable to connect to server"})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 10,
	}

	r2.Decr(database.Ctx, c.IP())

	val, _ = r.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}
