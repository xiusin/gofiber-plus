package wrapper

import (
	"errors"
	"log"
	"os"
	"reflect"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

type GroupRouter struct {
	NativeRouter fiber.Router
	wrapper      *AppWrapper
	Logger       LoggerAbstract
	CtrlName     string
}

func NewGroupRouter(router fiber.Router, wrapper *AppWrapper, Name string) *GroupRouter {
	return &GroupRouter{NativeRouter: router, CtrlName: Name, wrapper: wrapper, Logger: log.New(os.Stdout, "[ERR]", log.LstdFlags)}
}

func (g *GroupRouter) GetMethodWrapHandler(method string) fiber.Handler {
	typeOf := g.wrapper.GetControllerType(g.CtrlName)

	return func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if data := recover(); data != nil {
				err = ctx.JSON(fiber.Map{
					"status": fiber.StatusInternalServerError,
					"code":   fiber.StatusInternalServerError,
					"msg":    data.(error).Error(),
				})

				g.Logger.Print(string(debug.Stack()))
			}
		}()

		controller := reflect.New(typeOf)
		c := controller.Interface().(ControllerAbstract)
		c.SetCtx(ctx)
		c.Init()

		values := controller.MethodByName(method).Call(nil)

		if len(values) != 1 {
			panic(errors.New("请确定Handler有且只有一个返回值"))
		}

		if result := values[0].Interface(); result != nil {
			err = ctx.JSON(fiber.Map{
				"status": fiber.StatusInternalServerError,
				"code":   fiber.StatusInternalServerError,
				"msg":    err.Error(),
			})
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
