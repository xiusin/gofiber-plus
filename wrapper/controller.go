package wrapper

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

var ErrNeedImplements = errors.New(`please implement this method within the sub-structure`)

type Controller struct {
	*fiber.Ctx
}

func (b *Controller) Init() {}

func (b *Controller) SetCtx(ctx *fiber.Ctx) {
	b.Ctx = ctx
}

func (b *Controller) InitGroupRouter(router *GroupRouter) {
	panic(ErrNeedImplements)
}
