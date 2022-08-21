package wrapper

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type AppWrapper struct {
	*fiber.App

	reflectTypeData map[string]reflect.Type
}

func New(app *fiber.App) *AppWrapper {
	return &AppWrapper{app, map[string]reflect.Type{}}
}

func (wrapper *AppWrapper) GetControllerType(name string) reflect.Type {
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

	c.SetControllerName(name)
	c.SetRouterWrapper(wrapper)
	c.InitRouter(group)

	return wrapper
}
