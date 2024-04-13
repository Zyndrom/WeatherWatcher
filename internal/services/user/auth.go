package user

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserServices) RegisterUser(login, password string) error {
	exists, err := u.Storage.IsUserExist(login)
	if err != nil {
		log.Println(err)
		return err
	}
	if exists {
		return fmt.Errorf("user exists")
	}
	uuid := uuid.New()
	password, err = hashPassword(password)
	if err != nil {
		return fmt.Errorf("bcrypt error")
	}
	err = u.Storage.CreateNewUser(uuid.String(), login, password)
	if err != nil {
		return err
	}
	return nil
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UserServices) LoginUser(login, password string) (string, error) {
	uuid, hash, err := u.Storage.GetUserUUIDAndPassword(login)
	if err != nil {
		return "", err
	}
	if uuid == "" {
		return "", fmt.Errorf("wrong password or login")
	}
	if !checkPasswordHash(password, hash) {
		return "", fmt.Errorf("wrong password or login")
	}
	payload := jwt.MapClaims{
		"sub":   uuid,
		"login": login,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(u.jwtSecretKey))
	if err != nil {
		logrus.Debugf("%s", err.Error())
		return "", fmt.Errorf("jwt error")
	}
	return t, nil
}
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
