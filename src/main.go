package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/yishakk/fractage/src/config"
)

func main() {
	app := iris.New()
	port := strings.Trim(os.Getenv("PORT"), " ")
	if len(port) == 0 {
		port = "6060"
	}
	config.AddRoutes(app)
	app.Listen(fmt.Sprintf(":%s", port))
}
