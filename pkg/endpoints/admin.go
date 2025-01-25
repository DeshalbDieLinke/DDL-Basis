package endpoints

import (
	"ddl-server/pkg/database/models"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Returns JSON data for users in the database
func AdminPanel(c echo.Context) error {
	token, err := GetToken(c)
	if err != nil { 
		return c.JSON(401, map[string]string{"error": "No token provided"})
	}
	claims, err := GetTokenClaims(token)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid token")
	}
	if claims.AccessLevel != 0 {
		return c.String(http.StatusUnauthorized, "Insufficient access level. 0 Required "+fmt.Sprint(claims.AccessLevel)+" Provided")
	}
	db := c.Get("db").(*gorm.DB)

	var users []models.User
	db.Find(&users)
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Welcome to the admin panel", "users": users})
}
