package shared

import (
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId   string    `json:"userId"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.StandardClaims
}