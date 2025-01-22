package main

import (
	"deshalbdielinke/pkg/endpoints"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// db := database.StartDatabase()

	echo := echo.New()

	// Cache the certificate to prevent rate limiting
	echo.AutoTLSManager.Cache = autocert.DirCache("/root/certs")
	echo.AutoTLSManager.HostPolicy = autocert.HostWhitelist("api.deshalbdielinke.de")

	echo.GET("/", endpoints.HelloWorld)
	echo.GET("/content", endpoints.GetContent)
	echo.Logger.Fatal(echo.StartAutoTLS(":443"))

}
