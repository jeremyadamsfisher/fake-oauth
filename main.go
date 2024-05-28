package main

import (
	"github.com/jeremyadamsfisher/fake-user-cookie/cli"
	"github.com/jeremyadamsfisher/fake-user-cookie/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := cli.CommandLineInterface()
	app := echo.New()
	routes.AddToEcho(app, cfg)
	app.Logger.Fatal(app.Start(":4000"))
}
