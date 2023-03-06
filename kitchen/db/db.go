package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB(postgresUser, postgresPassword, postgresAddr, postgresPort, postgresDB string) (*DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", postgresAddr, postgresUser, postgresPassword, postgresDB, postgresPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Info("Connection to DB: ", dsn)

	return &DB{
		db: db,
	}, nil
}
