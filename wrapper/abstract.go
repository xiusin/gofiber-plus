package wrapper

import (
	"github.com/gofiber/fiber/v2"
)

type ControllerAbstract interface {
	SetCtx(*fiber.Ctx)
	InitGroupRouter(*GroupRouter)

	Init()
}

type LoggerAbstract interface {
	Print(v ...any)
}
