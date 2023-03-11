package test

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/xiusin/godi"
	"github.com/xiusin/gofiber-controller/wrapper"
)

type TestController struct {
	wrapper.Controller

	App *fiber.App `inject:"fiber.app"`
}

func (*TestController) Init() {}

func (*TestController) InitGroupRouter(w *wrapper.GroupRouter) {
	w.Get("/list", "List")
}

func (t *TestController) List() error {
	fmt.Println(t.App)
	return nil
}

func Test_App(t *testing.T) {
	app := fiber.New()

	godi.Instance("fiber.app", app)

	routerWrapper := wrapper.New(app)
	routerWrapper.WrapperHandler("/common", new(TestController))

	t.Log(app.Listen(":3001"))
}
