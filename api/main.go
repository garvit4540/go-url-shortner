package main

import (
	"github.com/garvit4540/go-url-shortner/routes"
	"github.com/garvit4540/go-url-shortner/trace"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"os"
)

func setupRoutes(app *fiber.App) {

	app.Get("/:url", routes.ResolveUrl)
	app.Post("/api/v1", routes.ShortenUrl)

}

func main() {

	err := godotenv.Load()
	if err != nil {
		trace.LogError(trace.ErrorLoadingEnvFiles, err, nil)
	}

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)

	// Start the app
	err = app.Listen(os.Getenv("APP_PORT"))
	if err != nil {
		trace.LogFatalError(trace.ErrorInStartingApp, err, nil)
	}
	trace.LogInfo(trace.AppStarted, map[string]interface{}{
		"port": os.Getenv("APP_PORT"),
	})
}
