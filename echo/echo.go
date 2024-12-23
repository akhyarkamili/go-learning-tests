package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Serve(address string) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(address))
}
