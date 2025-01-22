package endpoints

import (
	"ddl-server/pkg/utils"
	"encoding/json"
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
	materialJSON, err := json.Marshal(material)
	if err != nil {
		return c.String(http.StatusInternalServerError, "oopsie again")
	}
	return c.JSON(http.StatusOK, materialJSON)
}

func SearchContent(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func CreateContent(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
