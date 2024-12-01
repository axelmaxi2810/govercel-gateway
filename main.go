package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	route "gogateway/api/_routes"
)

func main() {
	// Create the Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Authorization Service - Local",
	})

	app.Use(healthcheck.New())
	// Initialize routes
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://*.micinproject.de, https://*.bandung.my.id",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			log.Printf(c.IP(), " Has Reach the limit")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  fiber.StatusTooManyRequests,
				"message": "TOO_MANY_REQUEST",
			})
		},
	}))

	route.RouteInit(app)

	port := ":4040" // Choose your desired local port
	log.Printf("Starting local server on http://localhost%s\n", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
