package main

import (
	"net/http"
	config "sentel/utility/config"

	"github.com/labstack/echo"
)

func main() {
	c := config.NewWithPath("")
	c.MustLoad(nil)
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello, world\n")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
