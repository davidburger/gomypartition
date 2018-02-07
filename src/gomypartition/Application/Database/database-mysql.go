package Database

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DbDriver struct {
	Host string
	User string
	Password string
	Port int
	Database string
	connection *sql.DB
}

func (db *DbDriver) GetConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Database)
}

func (db *DbDriver) GetConnection() (*sql.DB, error) {
	var err error
	if db.connection == nil {
		connection, err := sql.Open("mysql", db.GetConnectionString())

		if err == nil {
			db.connection = connection
		}
	}
	return db.connection, err
}

func (db *DbDriver) CloseConnection() {
	db.connection.Close()
}
