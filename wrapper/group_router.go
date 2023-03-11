package wrapper

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/xiusin/godi"
)

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
var InjectFailedFormat = "Injection service is causing an exception in the '%s' field. Reason: %s"

type GroupRouter struct {
	NativeRouter fiber.Router
	wrapper      *AppWrapper
	CtrlName     string
}

func NewGroupRouter(router fiber.Router, wrapper *AppWrapper, Name string) *GroupRouter {
	return &GroupRouter{NativeRouter: router, CtrlName: Name, wrapper: wrapper}
}

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

				g.wrapper.Logger.Print("Error message：", data,
					"\n ============ Stack ==============\n",
					string(msg[beginIndex+1:]),
					" =============  End  ==============\n")
			}
		}()

		controller := reflect.New(typeOf)
		c := controller.Interface().(ControllerAbstract)

		num := typeOf.NumField()
		for i := 0; i < num; i++ {
			field := typeOf.Field(i)
			name := field.Tag.Get("inject")
			if field.IsExported() && len(name) > 0 && godi.Exists(name) {
				valueOfField := controller.Elem().FieldByName(field.Name)
				if valueOfField.CanAddr() && valueOfField.IsNil() {
					func() {
						defer func() {
							if err := recover(); err != nil {
								g.wrapper.Logger.Print(fmt.Sprintf(InjectFailedFormat, field.Name, err))
							}
						}()
						valueOfField.Set(reflect.ValueOf(godi.MustGet(name)))
					}()
				}
			}
		}

		c.SetCtx(ctx)
		c.Init()

		methodRef := controller.MethodByName(method)
		if !methodRef.IsValid() {
			return ErrResponseHandler(ctx, `Method '`+method+`' does not exist.`)
		}
		if values := methodRef.Call(nil); len(values) != 1 {
			return ErrResponseHandler(ctx, `Please make sure that the handler has only one return value.`)
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
