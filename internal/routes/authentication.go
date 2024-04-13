package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserServices interface {
	RegisterUser(login, password string) error
	LoginUser(login, password string) (string, error)
}

func RegisterAuthentication(router *gin.Engine, usrService UserServices) {
	router.POST("/register", registration(usrService))
	router.POST("/login", login(usrService))
}

type userRequest struct {
	Login    string `json:"user"`
	Password string `json:"password"`
}

func registration(usrService UserServices) func(c *gin.Context) {
	return func(c *gin.Context) {
		var json userRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		err := usrService.RegisterUser(json.Login, json.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad login or password"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "user created"})
	}
}

func login(usrService UserServices) func(c *gin.Context) {
	return func(c *gin.Context) {
		var json userRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		jwt, err := usrService.LoginUser(json.Login, json.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wrong login or password"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Token": jwt,
		})
	}
}
