package wrapper

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var ErrFormat = errors.New("please pass in the correct format (method)")

type Controller struct {
	*fiber.Ctx
	wrapper        *AppWrapper
	ControllerName string
	MethodAction   int
}

func (b *Controller) Init() {}

func (b *Controller) InitRouter(router fiber.Router) {
	panic(errors.New("please implement in the sub controller"))
}

func (b *Controller) SetControllerName(name string) {
	b.ControllerName = name
}

func (b *Controller) SetRouterWrapper(wrapper interface{}) {
	b.wrapper = wrapper.(*AppWrapper)
}

func (b *Controller) warpHandler(method string) fiber.Handler {
	methods := strings.Split(method, "::")
	var name string

	if len(methods) == 1 {
		methods = append([]string{b.ControllerName}, methods[0])
	}
	if len(methods) == 2 {
		name = methods[0]
		if !strings.HasSuffix(name, "Controller") {
			name += "Controller"
		}
	} else {
		panic(ErrFormat)
	}

	typeOf := b.wrapper.GetControllerType(name)
	return func(ctx *fiber.Ctx) error {
		controller := reflect.New(typeOf)
		c := controller.Interface().(ControllerAbstract)
		c.SetCtx(ctx)
		c.SetRouterWrapper(b.wrapper)

		c.Init()

		values := controller.MethodByName(methods[1]).Call(nil)

		if result := values[0].Interface(); result != nil {
			return c.Ctx().JSON(fiber.Map{"status": fiber.StatusInternalServerError, "msg": result.(error).Error()})
		}

		return nil
	}
}
