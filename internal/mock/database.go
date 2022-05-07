// Package mock implement mocks
package mock

type Database struct {
	ConnectResp        error
	GetLastChangedResp string
}

func (d *Database) Connect(connectionString string) error {
	return d.ConnectResp
}

func (d *Database) GetLastChanged() string {
	return d.GetLastChangedResp
}

func (d *Database) CreateUser(id, date string) {
}
