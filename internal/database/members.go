package database

import (
	"github.com/joao3101/mailchimp-data-importer/internal/model"
)

func (d *database) GetLastChanged() string {
	var member model.Members
	d.Connector.Order("last_changed desc").First(&member)
	return member.LastChanged
}

func (d *database) CreateUser(id, date string) {
	member := model.Members{
		ID:          id,
		LastChanged: date,
	}
	d.Connector.Create(&member)
}
