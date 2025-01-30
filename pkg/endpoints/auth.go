package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ddl-server/pkg/database/models"
	"ddl-server/pkg/types"
	DDLErrors "ddl-server/pkg/types/errors"
	"ddl-server/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginJwt(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	secret := []byte(os.Getenv("JWT_SECRET"))
	var user models.User

	// Parse login request
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	if req.Email == "" || req.Password == "" {
		// Attempt to parse from Form
		req.Email = c.FormValue("email")
		req.Password = c.FormValue("password")
	}
	// token, err := GetToken(c)
	// if err != nil {
	// 	// No token provided - No Problem
	// }
	// claims, err := GetTokenClaims(token)
	// if err == nil {
	// 	if claims.Email == req.Email {
	// 		err := db.Where("email = ?", req.Email).First(&user).Error
	// 		if err != nil {
	// 			return c.JSON(401, map[string]string{"error": "Invalid "})

	// 		}
	// 		return c.JSON(http.StatusOK, map[string]string{"message": "Token valid for : " + claims.Email})
	// 	}
	// }

	// Validate credentials
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(401, map[string]string{"error": "Invalid email"})
		}
		return c.JSON(401, map[string]string{"error": "Invalid email or password"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password mismatch %v", user.Password+"| "+req.Password)
		return c.JSON(401, map[string]string{"error": "Invalid password"}) //TODO Add secure error message and checking
	}
	if true {
		// Generate JWT token
		claims := &types.JWTClaims{
			Email:       req.Email,
			AccessLevel: user.AccessLevel,
			ID:          user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // Token expires in 1 hour
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(secret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not generate token"})
		}
		// Return token to client
		// Cookify the token
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    signedToken,
			HttpOnly: true,
			Secure:   true,
			Expires:  claims.ExpiresAt.Time,
			SameSite: http.SameSiteNoneMode,
		})

		c.SetCookie(&http.Cookie{
			Name:     "id",
			Value:    fmt.Sprint(user.ID),
			HttpOnly: false,
			Secure:   true,
			Expires:  claims.ExpiresAt.Time,
			SameSite: http.SameSiteNoneMode,
		})

		return c.JSON(http.StatusOK, map[string]string{"message": "Login successful!", "email": req.Email})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
}

func Register(c echo.Context) error {
	type NewUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `json:"token"` // Token for the new user
	}
	var newUser NewUser

	// Default access level to user (3)
	accessLevel := 3

	err := c.Bind(&newUser)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	if newUser.Email == "" || newUser.Password == "" || newUser.Token == "" {
		return c.JSON(400, map[string]string{"error": "Invalid request: Empty email | password | token"})
	}

	// Parse Claims from token
	//TODO: Verify this works despite the
	claims, err := GetTokenClaims(newUser.Token)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid token"})
	}
	if claims != nil {
		if claims.Email != newUser.Email {
			return c.JSON(401, map[string]string{"error": "Email does not match token"})
		}
		log.Printf("Access level: %v", claims.AccessLevel)
		accessLevel = claims.AccessLevel
	}

	db := c.Get("db").(*gorm.DB)
	var user models.User

	if db.Where("email = ?", newUser.Email).First(&user).RowsAffected > 0 {
		return c.JSON(401, map[string]string{"error": "Email already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to hash password"})
	}
	db.Create(&models.User{Email: newUser.Email, Password: string(hashedPassword), AccessLevel: accessLevel})
	return c.JSON(201, map[string]string{"message": "User created successfully: " + newUser.Email + " " + string(accessLevel)})
}

func Profile(c echo.Context) error {
	token, err := GetToken(c)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "No token provided"})
	}
	claims, err := GetTokenClaims(token)
	if err != nil {
		log.Printf("Error verifying token: %v", err)
		return c.JSON(401, map[string]string{"error": "Invalid token"})
	}
	if claims.Email == "" {
		return c.JSON(400, map[string]string{"error": "Invalid token"})
	}
	return c.JSON(200, map[string]interface{}{
		"email":   claims.Email,
		"message": "Welcome to your profile!",
	})
}

// Logout Sets the token cookie to none - effectively logs out the user
func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Second),
		SameSite: http.SameSiteNoneMode,
	})
	return c.JSON(200, map[string]string{"message": "Logout successful"})
}

// Check Returns 200 if user is logged in
func Check(c echo.Context) error {
	// token := c.Request().Header.Get("Authorization")
	token, err := GetToken(c)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "No token provided"})
	}
	claims, err := GetTokenClaims(token)
	if err != nil {
		log.Printf("Error verifying token: %v", err)
		return c.JSON(401, map[string]string{"error": "Error verifying token"})

	}
	return c.JSON(200, map[string]string{"message": "Token valid until: " + claims.ExpiresAt.Time.GoString(), "accessLevel": fmt.Sprint(claims.AccessLevel), "email": claims.Email, "id": fmt.Sprint(claims.ID)})
}

// / Returns a token for a new user based on the input email and access level. Admin Level access is required.
func NewUserToken(c echo.Context) error {
	// Check if the user is an admin
	userToken, err := GetToken(c)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "No token provided"})
	}
	claims, err := GetTokenClaims(userToken)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid token"})
	}
	if claims.AccessLevel != 0 {
		return c.JSON(401, map[string]string{"error": "Insufficient access level. 0 Required " + fmt.Sprint(claims.AccessLevel) + " Provided"})
	}

	type NewUser struct {
		Email       string `json:"email"`
		AccessLevel int    `json:"accessLevel"`
	}

	newUser := new(NewUser)
	if err := c.Bind(&newUser); err != nil {
		log.Printf("Error binding new user: %v", err)
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	// Generate a token for the new user
	token, err := utils.GenerateToken(newUser.Email, newUser.AccessLevel)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(200, map[string]string{"token": token, "message": "Token generated successfully"})
}

// Utility functions

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
