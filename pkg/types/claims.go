package types

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Email       string `json:"email"`
	AccessLevel int
	jwt.RegisteredClaims
}