package wrapper

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

type GroupRouter struct {
	NativeRouter fiber.Router
	wrapper      *AppWrapper
	CtrlName     string
}

func NewGroupRouter(router fiber.Router, wrapper *AppWrapper, Name string) *GroupRouter {
	return &GroupRouter{NativeRouter: router, CtrlName: Name, wrapper: wrapper}
}

func (g *GroupRouter) parseController(method string) (string, string) {
	methods := strings.Split(method, "::")
	var name string

	if len(methods) == 1 {
		methods = append([]string{g.CtrlName}, methods[0])
	}
	if len(methods) == 2 {
		name = methods[0]
		if !strings.HasSuffix(name, "Controller") {
			name += "Controller"
		}
	} else {
		panic(ErrFormat)
	}

	return name, methods[1]
}

func (g *GroupRouter) GetMethodWrapHandler(methodSign string) fiber.Handler {
	ctrl, method := g.parseController(methodSign)

	typeOf := g.wrapper.GetControllerType(ctrl)

	return func(ctx *fiber.Ctx) error {
		controller := reflect.New(typeOf)
		c := controller.Interface().(ControllerAbstract)
		c.SetCtx(ctx)

		c.Init()

		values := controller.MethodByName(method).Call(nil)

		if result := values[0].Interface(); result != nil {
			return ctx.JSON(fiber.Map{"status": fiber.StatusInternalServerError, "msg": result.(error).Error()})
		}

		return nil
	}
}

func (g *GroupRouter) Post(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Post(path, mws...)
	return g
}

func (g *GroupRouter) Get(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Get(path, mws...)
	return g
}

func (g *GroupRouter) Head(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Head(path, mws...)
	return g
}

func (g *GroupRouter) Put(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Put(path, mws...)
	return g
}

func (g *GroupRouter) Delete(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Delete(path, mws...)
	return g
}

func (g *GroupRouter) Connect(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Connect(path, mws...)
	return g
}

func (g *GroupRouter) Options(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Options(path, mws...)
	return g
}

func (g *GroupRouter) Patch(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Patch(path, mws...)
	return g
}

func (g *GroupRouter) Trace(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.Trace(path, mws...)
	return g
}
