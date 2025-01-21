package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

func Logout(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

func CheckAuth(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

