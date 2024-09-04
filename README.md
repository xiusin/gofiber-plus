# gofiber-plus
Enable fiber to support controller mode



# demo
```go
package test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xiusin/gofiber-plus/wrapper"
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

```
