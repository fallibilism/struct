package config

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DBConnectionError = "failed to connect to db"

func NewDbConnection(c *PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", c.Host, c.Username, c.Password, c.DBName, c.Port, c.SslMode, c.TimeZone)
	ticker := time.NewTicker(5 * time.Second)

	var db *gorm.DB
	var err error
	for range ticker.C {
		if db != nil {
			break
		}

		fmt.Printf("running db connection again")

		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
			// PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %v", DBConnectionError, err)
	}

	// migrations

	// ping
	db.Exec("SELECT 1")
	fmt.Println("Connected to database successfully")
	return db, nil
}
