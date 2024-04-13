package middleware

import (
	"GoWeatherMap/internal/config"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CheckAuth(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Token") != "" {
			token, err := jwt.Parse(c.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("not authorized")
				}
				return []byte(config.JwtSecretKey), nil
			})
			if err != nil {
				c.Set("Auth", false)
				return
			}
			if !token.Valid {
				c.Set("Auth", false)
				return
			}
			c.Set("Auth", true)
			claims := token.Claims.(jwt.MapClaims)
			c.Set("uuid", claims["sub"])
		}
	}
}
