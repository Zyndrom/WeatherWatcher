package postgres

import (
	"GoWeatherMap/internal/config"
	"GoWeatherMap/internal/model"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	db *sql.DB
}

func NewStorage(cfg *config.Config) *Postgres {

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	err = db.Ping()
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	psg := &Postgres{db: db}
	psg.CreateTables()
	return psg
}

func (p *Postgres) CreateTables() {
	query := `CREATE TABLE IF NOT EXISTS Users (uuid VARCHAR(36) PRIMARY KEY, login VARCHAR(50), password VARCHAR(512));`
	_, err := p.db.Exec(query)
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	query = `CREATE TABLE IF NOT EXISTS WeatherLocations 
	(id serial primary key, name varchar(50), user_uuid varchar(36), ru_name varchar(50), country_code varchar(20), lat double precision, lon double precision);`
	_, err = p.db.Exec(query)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

}

func (p *Postgres) CreateNewUser(uuid string, login string, password string) error {
	query := `INSERT INTO Users (uuid,login, password) VALUES ($1, $2, $3)`
	_, err := p.db.Exec(query, uuid, login, password)
	if err != nil {
		logrus.Info(err)
	}
	return err
}

func (p *Postgres) GetUserUUIDAndPassword(login string) (string, string, error) {
	query := `SELECT uuid, password FROM Users WHERE login=$1;`
	row := p.db.QueryRow(query, login)
	uuid := ""
	password := ""
	err := row.Scan(&uuid, &password)
	if err != nil {
		return uuid, password, err
	}
	return uuid, password, nil
}

func (p *Postgres) IsUserExist(login string) (bool, error) {
	query := `SELECT COUNT(*) FROM Users WHERE login=$1`
	row := p.db.QueryRow(query, login)
	var count int
	err := row.Scan(&count)
	if err != nil {
		logrus.Fatalln(err)
		return true, err
	}
	return count > 0, nil
}

func (p *Postgres) AddLocationForUser(userUUID string, loc model.Location) error {
	query := `INSERT INTO WeatherLocations (name, user_uuid, ru_name, country_code, lat, lon ) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := p.db.Exec(query, loc.Name, userUUID, loc.RuName, loc.Country, loc.Lat, loc.Lon)
	if err != nil {
		logrus.Fatalln(err)
	}
	return err

}

func (p *Postgres) GetAllUserLocations(userUUID string) ([]model.Location, error) {
	locations := []model.Location{}
	query := `SELECT id,name, ru_name, country_code, lat, lon FROM WeatherLocations WHERE user_uuid=$1`
	rows, err := p.db.Query(query, userUUID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		l := model.Location{}
		err := rows.Scan(&l.Id, &l.Name, &l.RuName, &l.Country, &l.Lat, &l.Lon)
		if err != nil {
			log.Printf("Get All Locs err: %s", err.Error())
			return locations, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

func (p *Postgres) DeleteLocationForUser(userUUID string, loc model.Location) error {
	query := `DELETE FROM WeatherLocations WHERE user_uuid=$1 AND id=$2`
	_, err := p.db.Exec(query, userUUID, loc.Id)
	if err != nil {
		return err
	}
	return err
}
