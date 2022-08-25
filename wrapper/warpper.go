package wrapper

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type AppWrapper struct {
	fiber.Router

	reflectTypeData map[string]reflect.Type
}

func New(app fiber.Router) *AppWrapper {
	return &AppWrapper{app, map[string]reflect.Type{}}
}

func (wrapper *AppWrapper[T]) GetControllerType(name string) reflect.Type {
	reflectType, exists := wrapper.reflectTypeData[name]

	if !exists {
		panic(fmt.Errorf("controller %s does not exist", name))
	}

	return reflectType
}

func (wrapper *AppWrapper) WrapperHandler(prefix string, c ControllerAbstract, mws ...fiber.Handler) *AppWrapper {
	group := wrapper.Group(prefix, mws...)
	toc := reflect.TypeOf(c)
	name := toc.Elem().Name()
	wrapper.reflectTypeData[name] = toc.Elem()

	c.InitGroupRouter(NewGroupRouter(group, wrapper, name))
	return wrapper
}
