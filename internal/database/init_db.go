package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
)

var DB *sql.DB

// ConnDB initializes the database connection
func ConnDB(db_env config.DatabaseConfig) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db_env.Host,
		strconv.Itoa(db_env.Port),
		db_env.User,
		db_env.Password,
		db_env.Name,
		db_env.SSLmode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("cannot connect to DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("DB ping failed: %v", err)
	}

	log.Println("Connected to the database.")
	return nil
}
