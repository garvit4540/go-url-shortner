package routes

import (
	"github.com/garvit4540/go-url-shortner/database"
	"github.com/garvit4540/go-url-shortner/trace"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveUrl(ctx *fiber.Ctx) error {

	bodyUrl := ctx.Params("url")
	redisKey := "url:" + bodyUrl

	redisClient := database.CreateClient(0)
	defer redisClient.Close()

	value, err := redisClient.Get(database.Ctx, redisKey).Result()
	if err == redis.Nil {
		trace.LogError(trace.ErrorKeyNotFoundInRedis, err, nil)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found on database"})
	} else if err != nil {
		trace.LogError(trace.ErrorConnectingToRedis, err, nil)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to DB"})
	}

	// increment the counter
	rInr := database.CreateClient(1)
	defer rInr.Close()
	_ = rInr.Incr(database.Ctx, "counter:redirect_counter")

	trace.LogInfo(trace.RedirectedSuccessfully, map[string]interface{}{"url": bodyUrl, "redirect_url": value})
	return ctx.Redirect(value, fiber.StatusTemporaryRedirect)

}
