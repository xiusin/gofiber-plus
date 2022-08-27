package wrapper

import (
	"errors"
	"log"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type GroupRouter struct {
	NativeRouter fiber.Router
	wrapper      *AppWrapper
	Logger       LoggerAbstract
	CtrlName     string
}

var ControllerSuffix = "Controller"

func NewGroupRouter(router fiber.Router, wrapper *AppWrapper, Name string) *GroupRouter {
	return &GroupRouter{NativeRouter: router, CtrlName: Name, wrapper: wrapper, Logger: log.New(os.Stdout, "[ERR]", log.LstdFlags)}
}

func (g *GroupRouter) parseController(method string) (string, string) {
	methods := strings.Split(method, "::")
	var name string

	if len(methods) == 1 {
		methods = append([]string{g.CtrlName}, methods[0])
	}
	if len(methods) == 2 {
		name = methods[0]
		if !strings.HasSuffix(name, ControllerSuffix) {
			name += ControllerSuffix
		}
	} else {
		panic(ErrFormat)
	}

	return name, methods[1]
}

func (g *GroupRouter) GetMethodWrapHandler(methodSign string) fiber.Handler {
	ctrl, method := g.parseController(methodSign)
	typeOf := g.wrapper.GetControllerType(ctrl)

	return func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				err = recoverErr.(error)
				g.Logger.Print(string(debug.Stack()))
			}
		}()

		controller := reflect.New(typeOf)
		c := controller.Interface().(ControllerAbstract)
		c.SetCtx(ctx)
		c.Init()

		values := controller.MethodByName(method).Call(nil)
		length := len(values)

		if length > 1 {
			panic(errors.New("请确定Handler最多只有一个返回值"))
		} else if length == 1 {
			if result := values[0].Interface(); result != nil {
				switch result := result.(type) {
				case error:
					err = ctx.JSON(fiber.Map{
						"status": fiber.StatusInternalServerError,
						"msg":    result.Error(),
					})
				}
			}
		}

		return err
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
