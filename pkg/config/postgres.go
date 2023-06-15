import (
	"gorm.io/gorm"
	pg "gorm.io/driver/postgres"
)

const DBConnectionError = "failed to connect to db"
func NewDbConnection(c *PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=%s", c.User, c.Password, c.DBName, c.Port, c.SslMode)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	  }), &gorm.Config{})
	if err != nil {
		return nil, DBConnectionError
	}

	// migrations

	// ping
	return db, nil
}