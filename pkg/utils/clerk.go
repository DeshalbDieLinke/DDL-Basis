package utils

import (
	// clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	jwtClerk "github.com/clerk/clerk-sdk-go/v2/jwt"
	user "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/labstack/echo/v4"
)

var CLERK_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyduuv6zad48bi6frjYkY
sHWm4E660zi3xeAm/XXQ/FAhmOs0i5XPldLnrcjTAQuBYp1nEJmLe+IP+bAOehFD
sSRHz13KLGCQ4SzY+koFn35uIjGLc9buh6bF3qBSmXUrY/RHs9P81VCA1zLNBL1m
iawJuPZRlcAQSOvWkC6/QNDIZvClOIKU7pB2UUOZhkdHNj7nA4lhXLZwhwbzLe4e
yo/m0i4lf3dtdUgTAZi2gmT0BIb/Ez90EHIyMoof3zBPqwWmiIdQ/dhrOZKLw5Ph
mZk5tfaH3582gSfM4LrNs1nsAfy05mKxd17YU/1QOkcodZxaE45+ZJR2t5d88NHt
bQIDAQAB
-----END PUBLIC KEY-----` ;

func ClerkMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "authorization header missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
		}
		tokenStr := parts[1]

		claims, err := jwtClerk.Verify(c.Request().Context(), &jwtClerk.VerifyParams{
			Token: tokenStr,
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		usr, err := user.Get(c.Request().Context(), claims.Subject) 
		if err != nil { 
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		c.Set("user", usr)
		return next(c)
	}
}

func GetUserFromContext(c echo.Context) (*clerk.User, error) { 
	authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return nil, c.JSON(http.StatusUnauthorized, map[string]string{ "error": "authorization header missing", })
	}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return nil,  c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
		}
		tokenStr := parts[1]

		claims, err := jwtClerk.Verify(c.Request().Context(), &jwtClerk.VerifyParams{
			Token: tokenStr,
		})

		if err != nil {
			return nil, c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		usr, err := user.Get(c.Request().Context(), claims.Subject) 
		if err != nil { 
			return nil, c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return usr, nil
}

type Metadata struct {
	Role *string `json:"role"`
	Official *bool `json:"official"`
	Upload *bool `json:"upload"`
}
func GetUserRoleData(c echo.Context) (*Metadata, error) { 
	usr, err := GetUserFromContext(c)
	if err != nil {
		return nil, err
	}

	metadata := Metadata{}

	metadataBytes, err := usr.PrivateMetadata.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}