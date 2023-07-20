package routes

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rtsoy/url-shortener/helpers"
	"github.com/rtsoy/url-shortener/redis"
	"os"
	"strconv"
	"time"
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
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	// Parse request body to the struct
	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // 400
			"error": "Invalid JSON",
		})
	}

	// Rate limiting
	r1 := redis.CreateClient(1)
	defer r1.Close()

	_, err := r1.Get(redis.Ctx, c.IP()).Result()
	if errors.Is(err, redis.NotFound) {
		r1.Set(redis.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute)
	} else {
		rateStr, _ := r1.Get(redis.Ctx, c.IP()).Result()
		rate, _ := strconv.Atoi(rateStr)

		if rate < 1 {
			limit, _ := r1.TTL(redis.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":           "Rate limit exceeded",
				"rate_limit_rest": limit / time.Minute,
			})
		}
	}

	// Validate URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // 400
			"error": "Invalid URL",
		})
	}

	// Check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{ // 503
			"error": fmt.Sprintf("URLs containing `%s` cannon be shorten to avoid abuse", body.URL),
		})
	}

	// Enforce HTTPS, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	// Check if the user has provided custom short URL
	var id string
	if body.CustomShort == "" {
		id = uuid.NewString()
	} else {
		id = body.CustomShort
	}

	r := redis.CreateClient(0)
	defer r.Close()

	// Check if the user provided short is already in use
	value, err := r.Get(redis.Ctx, id).Result()
	if value != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom short URL is already in use",
		})
	}
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	if err := r.Set(redis.Ctx, id, body.URL, body.Expiry*time.Hour); err.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r1.Decr(redis.Ctx, c.IP())

	rateRemainingStr, _ := r1.Get(redis.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(rateRemainingStr)

	ttl, _ := r1.TTL(redis.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Minute

	return c.Status(fiber.StatusOK).JSON(resp)
}
