package route

import (
	"github.com/gofiber/fiber/v2"

	controller "gogateway/api/_controller"
)

func RouteInit(r *fiber.App) {
	healthCheck(r)
	routeGateway(r)
}

func healthCheck(r *fiber.App) {
	r.All("/live", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}

func routeGateway(r *fiber.App) {
	gatewayController := &controller.GatewayController{}
	r.All("/*", gatewayController.Index)
}
