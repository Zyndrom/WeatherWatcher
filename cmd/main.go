package main

import (
	"GoWeatherMap/internal/config"
	"GoWeatherMap/internal/routes"
	"GoWeatherMap/internal/services/user"
	"GoWeatherMap/internal/services/weather"
	"GoWeatherMap/internal/storage"
	"fmt"
	"path"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	setLogger()
	router := gin.New()
	router.Use(CORSMiddleware())
	config := config.ReadConfig()
	storage := storage.New(config)
	UserServices := user.NewUserService(storage, config)
	weatherServices := weather.NewWeatherController(config, storage)
	routes.RegisterAuthentication(router, UserServices)
	routes.RegisterLocation(router, config, weatherServices, UserServices)

	router.Run(":8080")
}

func setLogger() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)
	formatter := &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
