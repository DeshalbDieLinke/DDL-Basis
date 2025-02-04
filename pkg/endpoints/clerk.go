package endpoints

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/labstack/echo/v4"
)

func UserCreated(c echo.Context) error {
		ctx := c.Request().Context()
	
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}

		usr, err := user.Get(ctx, claims.Subject)
		if err != nil {
			panic(err)
		}
		if usr == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User does not exist"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User created" + usr.ID})
	}

