package user

import (
	"GoWeatherMap/internal/model"
)

func (u *UserServices) AddLocationForUser(uuid string, loc model.Location) error {
	err := u.Storage.AddLocationForUser(uuid, loc)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServices) DeleteLocationForUser(uuid string, loc model.Location) error {
	err := u.Storage.DeleteLocationForUser(uuid, loc)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServices) GetAllUserLocations(userUUID string) ([]model.Location, error) {
	locs, err := u.Storage.GetAllUserLocations(userUUID)
	return locs, err
}
