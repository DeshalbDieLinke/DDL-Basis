package utils

import (
	"ddl-server/pkg/database/models"
	"ddl-server/pkg/types"
	DDLErrors "ddl-server/pkg/types/errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Target struct { 
	ContentItem *models.Content
	User *models.User
}

// VerifyPermissions checks if the required permissions are met.
//
// Parameters:
// - required: The level of permissions required.
// - forSelf: A boolean indicating if the permissions are being verified for the user themselves. THis will allow
// Returns: A boolean indicating if the user has the required permissions.
func VerifyPermissions(required int, c echo.Context, target Target) bool {
	token, err := GetToken(c)
	if err != nil {
		return false
	}
	claims, err := GetTokenClaims(token)
	if err != nil {
		return false
	}
	// Validate the user exists and has the required permissions
	db := c.Get("db").(*gorm.DB)
	user := models.User{}
	if err := db.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		return false
	}
	if user.AccessLevel != claims.AccessLevel {
		log.Printf("WARNING!! Access level mismatch: %v != %v", user.AccessLevel, claims.AccessLevel)
		return false
	} else if user.Email != claims.Email {
		log.Printf("WARNING!! Email mismatch: %v != %v", user.Email, claims.Email)
		return false
	}


	if user.AccessLevel == 0 { 
		return true
	} 
	if target.ContentItem != nil && target.ContentItem.AuthorID == user.ID {
		return true
	}
	if target.User != nil && target.User.ID == user.ID {
		return true
	}
	return false
}


// / GetTokenFromRequest returns the token from the request header
func GetTokenFromRequest(c echo.Context) (string, error) {
	cookie, err := c.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}
	return "", err
}

func GetTokenClaims(t string) (*types.JWTClaims, error) {
	// Check if a token was provided
	token, err := jwt.ParseWithClaims(t, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		log.Printf("Error verifying token: %v", err)
		return nil, DDLErrors.InvalidToken
	}
	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: Claims could not be parsed")
	}
	return claims, nil
}

func GetToken(c echo.Context) (string, error) {
	var err error

	// Get the token from the request if it was not provided
	tokenStr, err := GetTokenFromRequest(c)
	if err != nil || tokenStr == "" {
		log.Printf("Error getting token: %v", err)
		return "", DDLErrors.NoTokenProvided
	}
	return tokenStr, nil

}
