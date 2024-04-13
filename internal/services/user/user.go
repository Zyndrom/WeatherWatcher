package user

import (
	"GoWeatherMap/internal/config"
	"GoWeatherMap/internal/model"
)

type UserServices struct {
	jwtSecretKey string
	Storage      storage
}

type storage interface {
	CreateNewUser(UUID, login, password string) error
	IsUserExist(UUID string) (bool, error)
	GetUserUUIDAndPassword(login string) (string, string, error)
	GetAllUserLocations(userUUID string) ([]model.Location, error)
	AddLocationForUser(userUUID string, loc model.Location) error
	DeleteLocationForUser(userUUID string, loc model.Location) error
}

func NewUserService(storage storage, config *config.Config) *UserServices {
	return &UserServices{jwtSecretKey: config.JwtSecretKey, Storage: storage}
}
