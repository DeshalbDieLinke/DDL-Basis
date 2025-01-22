package main

import (
	"deshalbdielinke/pkg/endpoints"

	"github.com/labstack/echo/v4"
)

func main() {
	// db := database.StartDatabase()

	echo := echo.New()
	echo.GET("/", endpoints.HelloWorld)
	echo.GET("/content", endpoints.GetContent)
	echo.Logger.Fatal(echo.StartAutoTLS(":443"))

}
