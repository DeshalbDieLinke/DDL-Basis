package endpoints

import (
	echo "github.com/labstack/echo/v4"
	"net/http"
)

func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
