package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ddl-server/pkg/database/models"
	"ddl-server/pkg/types"

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

	tokenStr := c.Request().Header.Get("Authorization")

	// Check if a token was provided
	if tokenStr != "" && tokenStr != "undefined" {
		token, err := jwt.ParseWithClaims(tokenStr, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token logging in"})
		}

		claims, ok := token.Claims.(*types.JWTClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token: Claims could not be parsed"})
		}
		if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
				return c.JSON(401, map[string]string{"error": "Token valid, Invalid email"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Token valid for : " + claims.Email})

	}

	// Validate credentials 
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound { 
            return c.JSON(401, map[string]string{"error": "Invalid email"})
        }
		return c.JSON(401, map[string]string{"error": "Invalid email or password"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        log.Printf("Password mismatch %v", user.Password + "| " + req.Password)
		return c.JSON(401, map[string]string{"error": "Invalid password"}) //TODO Add secure error message and checking 
	}
	if true {
		// Generate JWT token
		claims := &types.JWTClaims{
			Email: req.Email,
			AccessLevel: 3,
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
		return c.JSON(http.StatusOK, map[string]string{"token": signedToken, "message": "Login successful!", "email": req.Email})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
}

func Register(c echo.Context) error {
    type NewUser struct {
        Email string `json:"email"`
        Password string `json:"password"`
		Token string `json:"token"`
    }
    var newUser NewUser;

    err:= c.Bind(&newUser); if err != nil { 
        return c.JSON(400, map[string]string{"error": "Invalid request: No email or password provided"})
    }
    if newUser.Email == "" || newUser.Password == "" || newUser.Token == "" {
        return c.JSON(400, map[string]string{"error": "Invalid request: Empty email | password | token"})
    }

	// Parse Claims from token
	claims, err := ParseToken(newUser.Token)
	if err != nil { 
		return c.JSON(401, map[string]string{"error": "Invalid token"})
	}
	if claims.Email != newUser.Email { 
		return c.JSON(401, map[string]string{"error": "Email does not match token"})
	}
	accessLevel := claims.AccessLevel

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
	user := c.Get("user").(*jwt.Token)       // Get the JWT token from context
	claims := user.Claims.(*types.JWTClaims)      // Extract claims
	if claims.Email == "" {
		return c.JSON(400, map[string]string{"error": "Invalid token"})
	}
	return c.JSON(200, map[string]interface{}{
		"email":   claims.Email,
		"message": "Welcome to your profile!",
	
	})
}

func ParseToken(tokenString string) (*types.JWTClaims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JST_TOKEN")), nil
	})
	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*types.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Returns 200 if user is logged in
func Check(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	claims, err := VerifyToken(string(token), c)
	if err != nil {
		log.Printf("Error verifying token: %v", err)
		return c.JSON(401, map[string]string{"error": "Invalid token"})

	}
	return c.JSON(200, map[string]string{"message": "Token valid until: " + claims.ExpiresAt.Time.GoString(), "accessLevel": fmt.Sprint(claims.AccessLevel), "email": claims.Email})
}


func VerifyToken(tokenStr string, c echo.Context) (*types.JWTClaims, error) {
		// Check if a token was provided
		if tokenStr != "" {
			token, err := jwt.ParseWithClaims(tokenStr, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is correct
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SWT_TOKEN")), nil
			})
	
			if err != nil || !token.Valid {
				log.Printf("Error verifying token: %v", err)
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token could not be verified")
			}
			
			claims, ok := token.Claims.(*types.JWTClaims)
			if !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized,"Invalid token: Claims could not be parsed")
				}
			return claims, nil
		}
		return nil, echo.NewHTTPError(http.StatusUnauthorized,"No token provided")
}