package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewDatabase(log *logrus.Logger) *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connectionStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		dbUser,
		dbPassword,
		dbName,
		dbHost,
		dbPort,
		dbSSLMode,
	)

	db, err := sqlx.Connect("postgres", connectionStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Verify the connection to the database is still alive
	err = db.Ping()
	if err != nil {
		panic("Failed to ping the database: " + err.Error())
	}

	log.Info("Successfully connected to the database! 🐘")

	return db
}
