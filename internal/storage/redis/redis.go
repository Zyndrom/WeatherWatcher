package redis

import (
	"GoWeatherMap/internal/config"
	"GoWeatherMap/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisBd struct {
	client *redis.Client
}

const (
	expirationTime = time.Hour * 24
)

func New(cfg *config.Config) *RedisBd {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	status := client.Ping(context.Background())
	if status.Err() != nil {
		logrus.Fatalf(status.Err().Error())
	}
	return &RedisBd{client: client}
}

func (r *RedisBd) GetLocationByName(name string) ([]model.Location, error) {
	locs := []model.Location{}
	val, err := r.client.Get(context.Background(), fmt.Sprintf("location_req:%s", name)).Result()
	if err != nil {
		return locs, err
	}
	err = json.Unmarshal([]byte(val), &locs)
	if err != nil {
		return locs, err
	}
	return locs, nil
}

func (r *RedisBd) SetLocationByName(name string, locs []model.Location) error {
	err := r.client.Set(context.Background(), fmt.Sprintf("location_req:%s", name), locs, expirationTime).Err()
	return err
}
