package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
)

var DB *sql.DB

func ConnDB(db_env config.DatabaseConfig) {
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
		log.Fatal("Cannot connect to DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	fmt.Println("Connected to the database.")

	sqlFilePath := "/home/zihad/coding/Developer-Assignment/migrations/000001_init_schema/up.sql"

	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("DB execute error: %v", err)
	}

	log.Println("Hurray, database schema created successfully.")
}
