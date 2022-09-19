package wrapper

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/xiusin/godi"
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

// ErrResponseHandler 错误处理函数
var ErrResponseHandler = func(ctx *fiber.Ctx, data string) error {
	return ctx.JSON(fiber.Map{
		"status": fiber.StatusInternalServerError,
		"code":   fiber.StatusInternalServerError,
		"msg":    data,
	})
}

var beginStackSubBytes = []byte("src/runtime/panic.go:")
var endStackSubBytes = []byte("github.com/gofiber/fiber/v2.(*App).next(")

func (g *GroupRouter) GetMethodWrapHandler(method string) fiber.Handler {
	typeOf := g.wrapper.GetControllerType(g.CtrlName)
	return func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if data := recover(); data != nil {
				err = ErrResponseHandler(ctx, fmt.Sprintf("%s", data))

				stack := debug.Stack()
				beginIndex, endIndex := bytes.Index(stack, beginStackSubBytes), bytes.Index(stack, endStackSubBytes)
				var msg = stack[beginIndex:endIndex]

				stack = nil
				beginIndex = bytes.Index(msg, []byte{'\n'})

				g.Logger.Print("错误信息：", data,
					"\n ============ 堆栈 ==============\n",
					string(msg[beginIndex+1:]),
					" ============= 结束 ==============\n")
			}
		}()

		controller := reflect.New(typeOf)

		valueOf := reflect.ValueOf(controller)

		c := controller.Interface().(ControllerAbstract)

		fieldNum := typeOf.NumField()
		for i := 0; i < fieldNum; i++ {
			field := typeOf.Field(i)
			valueOfField := valueOf.Field(i)
			if !valueOfField.CanAddr() || !valueOfField.IsNil() {
				continue
			}
			if serviceName := field.Tag.Get("di"); len(serviceName) > 0 && godi.Exists(serviceName) {
				func() {
					defer func() {
						if err := recover(); err != nil {
							g.Logger.Print("reflect failed", err)
						}
					}()
					valueOfField.Set(reflect.ValueOf(godi.MustGet(serviceName)))
				}()
			}
		}

		c.SetCtx(ctx)

		c.Init()

		if values := controller.MethodByName(method).Call(nil); len(values) != 1 {
			panic(errors.New("请确定Handler有且只有一个返回值"))
		} else if result := values[0].Interface(); result != nil {
			return ErrResponseHandler(ctx, fmt.Sprintf("%s", result))
		}
		return nil
	}
}

func (g *GroupRouter) All(path string, method string, mws ...fiber.Handler) *GroupRouter {
	mws = append([]fiber.Handler{g.GetMethodWrapHandler(method)}, mws...)
	g.NativeRouter.All(path, mws...)
	return g
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
