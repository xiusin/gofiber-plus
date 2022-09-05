package wrapper

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var ErrFormat = errors.New("please pass in the correct format (method)")

type Controller struct {
	*fiber.Ctx
}

func (b *Controller) Init() {}

func (b *Controller) SetCtx(ctx *fiber.Ctx) {
	b.Ctx = ctx
}

func (b *Controller) InitGroupRouter(router *GroupRouter) {
	panic(errors.New("please implement in the sub controller"))
}
