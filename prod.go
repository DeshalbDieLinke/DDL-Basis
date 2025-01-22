package main

import (
	"ddl-server/pkg/endpoints"

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
	echo.Logger.Fatal(echo.StartTLS(":8080", "/etc/letsencrypt/live/api.deshalbdielinke.de/fullchain.pem", "/etc/letsencrypt/live/api.deshalbdielinke.de/privkey.pem"))

}
