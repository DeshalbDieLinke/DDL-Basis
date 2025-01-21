package main

import (
	"deshalbdielinke/pkg/endpoints"

	"github.com/labstack/echo/v4"
)

func main() {
	// db := database.StartDatabase()
	
	
	e := echo.New()
	e.GET("/", endpoints.HelloWorld)
	e.Logger.Fatal(e.StartAutoTLS(":443"))

}
