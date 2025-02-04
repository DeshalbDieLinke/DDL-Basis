package endpoints

import (
	"ddl-server/pkg/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUsers Returns JSON data for users in the database
func GetUsers(c echo.Context) error {
	token, err := utils.GetToken(c)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "No token provided"})
	}
	claims, err := utils.GetTokenClaims(token)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid token")
	}
	if claims.AccessLevel != 0 {
		return c.String(http.StatusUnauthorized, "Insufficient access level. 0 Required "+fmt.Sprint(claims.AccessLevel)+" Provided")
	}
	// db := c.Get("db").(*gorm.DB)

	// var users []models.User
	// db.Find(&users)
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Welcome to the admin panel", "users": "NOT IMPLEMENTED"})
}
