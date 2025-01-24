package endpoints

import (
	"ddl-server/pkg/utils"
	"io"
	"net/http"
	"os"

	echo "github.com/labstack/echo/v4"
)



func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

//TODO Add support for query parameters and search as well as returning the content from the database
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
	// Get the database connection
	// db := c.Get("db").(*gorm.DB)

	// Get the FormData
	// title := c.FormValue("title")
	// description := c.FormValue("description")
	// topics := c.FormValue("topics")
	// official := c.FormValue("official")
	// Get the file 
	image, err  := c.FormFile("image")
	if err != nil { 
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "No file provided"})
	}

	// Save the file
	src, err := image.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error opening file"})
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(image.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating file"})
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}


	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
}
