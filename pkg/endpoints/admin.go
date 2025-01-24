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
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.String(http.StatusUnauthorized, "No token provided")
	}
	claims, err := ParseToken(token)
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
