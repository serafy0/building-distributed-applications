package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct{}
type Claims struct {
	Username string `json:"username"`
	jwt.Claims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {}
