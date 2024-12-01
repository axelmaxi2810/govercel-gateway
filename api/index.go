package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	route "gogateway/api/_routes"
)

// Handler is the main entry point for Vercel.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set the request path in `*fiber.Ctx`
	r.RequestURI = r.URL.String()
	handler().ServeHTTP(w, r)
}

// handler initializes the Fiber application and sets up routes
func handler() http.HandlerFunc {
	// Database Initialization
	//	database.DatabaseInitSSL()

	// Database Migration
	// migration.RunMigrations()

	// Create the Fiber application
	app := fiber.New(fiber.Config{
		AppName: "Gateway - Service",
	})

	// Panic Handler
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://*.micinproject.de, https://*.bandung.my.id, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(healthcheck.New())
	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.Get("Cf-Connecting-Ip")
			if ip == "" {
				ip = c.IP()
			}
			return ip
		},
		LimitReached: func(c *fiber.Ctx) error {
			log.Printf(c.IP(), " Has Reach the limit")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  fiber.StatusTooManyRequests,
				"message": "TOO_MANY_REQUEST",
			})
		},
	}))

	// Initialize routes
	route.RouteInit(app)

	return adaptor.FiberApp(app)
}
