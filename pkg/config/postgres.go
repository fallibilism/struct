package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DBConnectionError = "failed to connect to db"

func NewDbConnection(c *PostgresConfig) (*gorm.DB, error) {
	DSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", c.Username, c.Password, c.Host, c.Port, c.DBName, c.SslMode)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}))

	// DSN: "postgresql://doadmin:AVNS_-OK3KDjBah18nx3cALr@db-postgresql-fra1-42722-do-user-9369539-0.b.db.ondigitalocean.com:25060/defaultdb?sslmode=require",
	// DriverName: "cloudsqlpostgres",
	// Conn:       sqlDb,
	// PreferSimpleProtocol: true, // disables implicit prepared statement usage
	// defer db.Close()

	if err != nil {
		return nil, fmt.Errorf("%s: %v", DBConnectionError, err)
	}

	// migrations

	// ping
	db.Exec("SELECT 1")
	fmt.Println("Connected to database successfully")
	return db, nil
}
