package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rtsoy/url-shortener/redis"
)

func ResolveURL(c *fiber.Ctx) error {
	// Get the shortened URL from the URL
	url := c.Params("url")

	r := redis.CreateClient(0)

	// Query to find original URL
	value, err := r.Get(redis.Ctx, url).Result()
	if errors.Is(err, redis.NotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ // 404
			"error": "URL not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ // 500
			"error": "Internal Server Error",
		})
	}

	// Increment the counter
	toIncr := redis.CreateClient(1)
	defer toIncr.Close()

	_ = toIncr.Incr(redis.Ctx, "counter")

	// Redirect to original URL
	return c.Redirect(value, fiber.StatusMovedPermanently) // 301
}
