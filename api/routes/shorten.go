package routes

import (
	govalidator "github.com/asaskevich/govalidator"
	"github.com/garvit4540/go-url-shortner/database"
	"github.com/garvit4540/go-url-shortner/helpers"
	"github.com/garvit4540/go-url-shortner/trace"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"os"
	"strconv"
	"time"
)

func ShortenUrl(ctx *fiber.Ctx) error {

	// Parse Body
	body := new(request)
	if err := ctx.BodyParser(&body); err != nil {
		trace.LogError(trace.ErrorParsingBodyToJSON, err, nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Implement Rate Limiting
	redisClient := database.CreateClient(1)
	defer redisClient.Close()
	redisKey := "user:" + ctx.IP()
	val, err := redisClient.Get(database.Ctx, redisKey).Result()
	if err == redis.Nil {
		trace.LogError(trace.ErrorKeyNotFoundInRedis, err, nil)
		_ = redisClient.Set(database.Ctx, redisKey, os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else if err != nil {
		trace.LogError(trace.ErrorConnectingToRedis, err, nil)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to DB"})
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := redisClient.TTL(database.Ctx, redisKey).Result()
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(map[string]interface{}{
				"error":            "Rate Limit Exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// Check if the input is an actual url or not
	if govalidator.IsURL(body.URL) == false {
		trace.LogError(trace.ErrorInvalidUrl, err, nil)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "url provided is not valid"})
	}

	// Check for any domain error
	if helpers.RemoveDomainError(body.URL) == false {
		trace.LogError(trace.ErrorSelfDomainLoopPrevented, err, nil)
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "haha .. nice try"})
	}

	// Implementing the Url Shortening Functionality
	id, err := shortenUrlHelper(ctx, body)
	if err != nil {
		return err
	}

	// Enforce http // ssl
	body.URL = helpers.EnforceHttp(body.URL)

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}
	redisKey = "user:" + ctx.IP()
	redisClient.Decr(database.Ctx, redisKey)
	val, _ = redisClient.Get(database.Ctx, redisKey).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := redisClient.TTL(database.Ctx, redisKey).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return nil
}

func shortenUrlHelper(ctx *fiber.Ctx, body *request) (string, error) {

	var id string
	if body.CustomShort != "" {
		id = body.CustomShort
	} else {
		id = uuid.New().String()[:6]
	}

	redisClient := database.CreateClient(0)
	defer redisClient.Close()
	redisKey := "url:" + id

	val, err := redisClient.Get(database.Ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		trace.LogError(trace.ErrorConnectingToRedis, err, nil)
		return "", ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to DB"})
	} else if val != "" {
		trace.LogError(trace.ErrorCustomShortAlreadyExists, err, nil)
		return "", ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL short already in use"})
	}

	if body.Expiry == 0 {
		body.Expiry = 24 // default expiry of 24 hours
	}

	err = redisClient.Set(database.Ctx, redisKey, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		trace.LogError(trace.ErrorConnectingToRedis, err, nil)
		return "", ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to connect to server"})
	}

	return id, nil

}
