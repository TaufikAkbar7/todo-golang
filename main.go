package main

import (
	"fmt"
	"golang-todo/internal/config"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using OS environment variables")
	}
	logrus := config.NewLogger()
	db := config.NewDatabase(logrus)
	validate := config.NewValidator()
	app := config.NewFiber()

	defer db.Close()

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      logrus,
		Validate: validate,
	})
	port := os.Getenv("APP_PORT")
	errApp := app.Listen(fmt.Sprintf(":%s", port))
	if errApp != nil {
		log.Fatalf("Failed to start server: %v", errApp)
	}
}
