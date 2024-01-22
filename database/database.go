package database

import (
	"database/sql"
	"fmt"
	"log"

	"bookings.com/m/config"
	_ "github.com/lib/pq"
)

func InitialiseConnection(l *log.Logger) (*sql.DB, error) {
	configFileName := "config/config.json"

	config, err := config.ReadConfigFile(configFileName)
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// test db is connected using ping
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	l.Println("Database Connection Successful!")
	return db, nil
}
