package main

import (
	"ddl-server/pkg/database"
	"ddl-server/pkg/database/models"
	"ddl-server/pkg/endpoints"
	"ddl-server/pkg/types"
	"ddl-server/pkg/utils"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	jwtE "github.com/labstack/echo-jwt/v4"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)


func main() {
	err:= godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	log.Print("Starting server")

	db := database.StartDatabase()

	e := echo.New()

	var SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))

	// Register middleware
	e.Use(database.GormMiddleware(db))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://deshalbdielinke.de", "https://api.deshalbdielinke.de"},
		AllowMethods: []string{"GET", "POST"},
	}))
	e.Use(jwtE.WithConfig(jwtE.Config{
		Skipper: func(c echo.Context) bool {
			return !strings.Contains(c.Path(), "/auth"); 
		},
		SigningKey: SECRET_KEY,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &types.JWTClaims{}
		},
		TokenLookup: "header:Authorization",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			return jwt.ParseWithClaims(auth, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Ensure the signing method is correct
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return SECRET_KEY, nil
			})
		},
	}));

	// Register endpoints
	e.POST("/register", endpoints.Register)
	e.POST("/login", endpoints.LoginJwt)
	e.GET("/", endpoints.HelloWorld)
	e.GET("/content", endpoints.GetContent)
	// Register a catch-all route
	e.Any("/*", func(context echo.Context) error { 
		return context.String(404, "Not Found")
	})


	// Register Restricted Endpoints
	restricted := e.Group("/auth")
	restricted.GET("/profile", endpoints.Profile)
	restricted.POST("/create", endpoints.CreateContent)
	restricted.GET("/admin", endpoints.SearchContent)
	restricted.GET("/check", endpoints.Check)
	restricted.GET("/*", func(c echo.Context) error {
		log.Printf("Authenticated request")
		return c.JSON(200, map[string]string{"message": "Authenticated request"})	
	})
	
	// Check if DB is empty and create a default user
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 { 
		adminEmail := os.Getenv("INIT_EMAIL")
		if adminEmail == "" {
			log.Fatal("No admin email provided")
		}
		initToken, _ := utils.GenerateToken(adminEmail, 0)
		log.Printf("Creating default user for %v: %s",  adminEmail, initToken)
	} else {
		log.Printf("Count: %d", count)
	}

	// Cache the certificate to prevent rate limiting
	e.AutoTLSManager.Cache = autocert.DirCache("/root/certs")
	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("api.deshalbdielinke.de")
	
	e.Logger.Fatal(e.StartTLS(":8080", "/etc/letsencrypt/live/api.deshalbdielinke.de/fullchain.pem", "/etc/letsencrypt/live/api.deshalbdielinke.de/privkey.pem"))
	// e.Logger.Fatal(e.Start(":8080"))

}
