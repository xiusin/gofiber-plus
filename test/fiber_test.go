package test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xiusin/gofiber-controller/wrapper"
	"testing"
)

type TestController struct {
	wrapper.Controller
}

func (*TestController) InitGroupRouter(*wrapper.GroupRouter) {

}

func Test_App(t *testing.T) {

	app := fiber.New()

	routerWrapper := wrapper.New(app)
	routerWrapper.WrapperHandler("/common", new(TestController))
}
