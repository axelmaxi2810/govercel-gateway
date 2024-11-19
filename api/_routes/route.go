package route

import (
	"github.com/gofiber/fiber/v2"

	controller "gogateway/api/_controller"
)

func RouteInit(r *fiber.App) {
	routeGateway(r)
}

func routeGateway(r *fiber.App) {
	gatewayController := &controller.GatewayController{}
	r.All("/*", gatewayController.Index)
}
