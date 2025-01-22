package main

import (
	"ddl-server/pkg/endpoints"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// db := database.StartDatabase()

	e := echo.New()

	// Cache the certificate to prevent rate limiting
	e.AutoTLSManager.Cache = autocert.DirCache("/root/certs")
	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("api.deshalbdielinke.de")

	e.GET("/", endpoints.HelloWorld)
	e.GET("/content", endpoints.GetContent)
	e.Any("/*", func(context echo.Context) error { 
		return context.String(404, "Not Found")
	})
	e.Logger.Fatal(e.StartTLS(":8080", "/etc/letsencrypt/live/api.deshalbdielinke.de/fullchain.pem", "/etc/letsencrypt/live/api.deshalbdielinke.de/privkey.pem"))

}
