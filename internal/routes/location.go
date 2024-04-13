package routes

import (
	"GoWeatherMap/internal/config"
	"GoWeatherMap/internal/middleware"
	"GoWeatherMap/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationRouter struct {
}

type userServices interface {
	AddLocationForUser(uuidUser string, loc model.Location) error
	DeleteLocationForUser(uuidUser string, loc model.Location) error
	GetAllUserLocations(userUUID string) ([]model.Location, error)
}
type weatherServices interface {
	FindLocationByName(name string) ([]model.Location, error)
	CurrentWeatherInLocation(loc model.Location) (model.Weather, error)
}

func RegisterLocation(router *gin.Engine, cfg *config.Config, weather weatherServices, userLocController userServices) {
	loc := &LocationRouter{}
	locGroup := router.Group("/location", middleware.CheckAuth(cfg))
	locGroup.GET("/find", loc.findLocation(weather))
	locGroup.GET("/all", loc.getAllUserLocation(userLocController))
	locGroup.POST("/add", loc.addWeatherLocation(userLocController))
	locGroup.DELETE("/delete", loc.removeWeatherLocation(userLocController))
	locGroup.POST("/current-weather", loc.locationWeather(weather))
}

// Send model.Location list 0-3
func (l *LocationRouter) findLocation(weatherServices weatherServices) func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.GetBool("Auth")
		if !auth {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		name := c.DefaultQuery("name", "Moscow")
		locs, err := weatherServices.FindLocationByName(name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		c.JSON(http.StatusOK, locs)
	}
}

func (l *LocationRouter) addWeatherLocation(userLocController userServices) func(c *gin.Context) {
	return func(c *gin.Context) {
		if ok, exists := c.Get("Auth"); !exists || !ok.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		uuid, exists := c.Get("uuid")
		if !exists || uuid == nil || uuid.(string) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var json model.Location
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		userLocController.AddLocationForUser(uuid.(string), json)
	}
}
func (l *LocationRouter) getAllUserLocation(userLocController userServices) func(c *gin.Context) {
	return func(c *gin.Context) {
		if ok, exists := c.Get("Auth"); !exists || !ok.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		uuid, exists := c.Get("uuid")
		if !exists || uuid == nil || uuid.(string) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		locs, err := userLocController.GetAllUserLocations(uuid.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		c.JSON(http.StatusAccepted, locs)
	}
}

func (l *LocationRouter) removeWeatherLocation(userLocController userServices) func(c *gin.Context) {
	return func(c *gin.Context) {
		if ok, exists := c.Get("Auth"); !exists || !ok.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		uuid, exists := c.Get("uuid")
		if !exists || uuid == nil || uuid.(string) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var json model.Location
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		userLocController.DeleteLocationForUser(uuid.(string), json)
	}
}

func (l *LocationRouter) locationWeather(weatherServices weatherServices) func(c *gin.Context) {
	return func(c *gin.Context) {
		if ok, exists := c.Get("Auth"); !exists || !ok.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		uuid, exists := c.Get("uuid")
		if !exists || uuid.(string) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var loc model.Location
		if err := c.ShouldBindJSON(&loc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		weather, err := weatherServices.CurrentWeatherInLocation(loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		c.JSON(http.StatusAccepted, weather)
	}
}
