package types

import (
	"net/http"

	e "github.com/labstack/echo/v4"
) 

var (
	InvalidToken  = e.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	EmailAlreadyExists = e.NewHTTPError(http.StatusConflict, "Email already exists")
	InvalidRequest     = e.NewHTTPError(http.StatusBadRequest, "Invalid request")
)