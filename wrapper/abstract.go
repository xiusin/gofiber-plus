package wrapper

import (
	"github.com/gofiber/fiber/v2"
)

type ControllerAbstract interface {
	InitRouter(fiber.Router)
	SetRouterWrapper(interface{})
	SetControllerName(string)

	Init()
}
