package endpoints

import (
	"log"

	"github.com/labstack/echo/v4"
)

func UserCreated(c echo.Context) error {
		ctx := c.Request().Context()
	
		log.Fatalf("User created" + ctx.Value("type").(string))
		// if usr == nil {
		// 	return c.JSON(http.StatusNotFound, map[string]string{"message": "User does not exist"})
		// }

		// return c.JSON(http.StatusOK, map[string]string{"message": "User created" + usr.ID})
		return c.JSON(200, map[string]string{"message": "User created"})
	}

