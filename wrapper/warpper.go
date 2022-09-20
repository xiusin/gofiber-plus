package wrapper

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type AppWrapper struct {
	fiber.Router
	reflectTypeData map[string]reflect.Type
	Logger          LoggerAbstract
}

func New(app fiber.Router) *AppWrapper {
	return &AppWrapper{app, map[string]reflect.Type{}, log.New(os.Stdout, "[ERR]", log.LstdFlags)}
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
	c.InitGroupRouter(NewGroupRouter(group, wrapper, name))
	return wrapper
}
