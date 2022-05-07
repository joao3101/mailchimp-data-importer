package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database interface {
	Connect(connectionString string) error
	GetLastChanged() string
	CreateUser(id, date string)
}

type database struct {
	Connector *gorm.DB
}

func NewDatabase(connectionString string) (Database, error) {
	Connector, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &database{
		Connector: Connector,
	}, nil
}

//Connect creates MySQL connection
func (d *database) Connect(connectionString string) error {
	var err error
	d.Connector, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Println("Connection was successful!!")
	return nil
}
