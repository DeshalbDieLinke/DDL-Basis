package endpoints

import (
	"ddl-server/pkg/utils"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func GetContent(c echo.Context) error {
	material, err := utils.GetMaterial()
	if err != nil {
		return c.String(http.StatusInternalServerError, "oopsie")
	}
	return c.JSON(http.StatusOK, material)
}

func SearchContent(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func CreateContent(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
