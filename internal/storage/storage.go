package storage

import (
	"GoWeatherMap/internal/config"
	"GoWeatherMap/internal/model"
	"GoWeatherMap/internal/storage/postgres"
	"GoWeatherMap/internal/storage/redis"
)

type Storage struct {
	postgres *postgres.Postgres
	redis    *redis.RedisBd
}

func New(cfg *config.Config) *Storage {
	storage := &Storage{
		postgres: postgres.NewStorage(cfg),
		redis:    redis.New(cfg),
	}
	return storage
}

func (s *Storage) GetAllUserLocations(userUUID string) ([]model.Location, error) {
	return s.postgres.GetAllUserLocations(userUUID)
}
func (s *Storage) AddLocationForUser(userUUID string, loc model.Location) error {
	return s.postgres.AddLocationForUser(userUUID, loc)
}
func (s *Storage) IsUserExist(login string) (bool, error) {
	return s.postgres.IsUserExist(login)
}
func (s *Storage) GetUserUUIDAndPassword(login string) (string, string, error) {
	return s.postgres.GetUserUUIDAndPassword(login)
}

func (s *Storage) CreateNewUser(UUID, login, password string) error {
	return s.postgres.CreateNewUser(UUID, login, password)
}

func (s *Storage) DeleteLocationForUser(userUUID string, loc model.Location) error {
	return s.postgres.DeleteLocationForUser(userUUID, loc)
}

func (s *Storage) GetLocationByName(name string) ([]model.Location, error) {
	return s.redis.GetLocationByName(name)
}
func (s *Storage) SetLocationByName(name string, locs []model.Location) error {
	return s.redis.SetLocationByName(name, locs)
}
