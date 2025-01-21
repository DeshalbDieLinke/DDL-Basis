package main

import (
	"deshalbdielinke/pkg/database"
	"deshalbdielinke/pkg/endpoints"

	"github.com/labstack/echo/v4"
)

func main() {
	database.StartDatabase()
	e := echo.New()
	e.GET("/", endpoints.HelloWorld)
	e.Logger.Fatal(e.Start(":1312"))

}
