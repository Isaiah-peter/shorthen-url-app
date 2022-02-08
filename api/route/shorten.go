package route

import (
	"time"

	"shorten-url-with-redis/helper"

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
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse json"})
	}
	//implement rate limiting

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

}
