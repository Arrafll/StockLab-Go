package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Arrafll/StockLab-Go/internal/config"
	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/routes"
	_ "github.com/Arrafll/StockLab-Go/internal/routes"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

// @title StockLab
// @version 1.0
// @description Dokumentasi API StockLab
// @contact.name Andre R
// @contact.email andrerafli83@gmail.com
// @license.name dreowsy
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" JWT
func main() {
	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	cfg := config.Load()

	Info.Println("Connecting to database...")

	_, err = db.Connect(cfg)
	if err != nil {
		Error.Println("Failed to connect to database: %v", err)
		os.Exit(1)
	}

	Info.Println("Connected to database...")
	route := routes.RegisterRoutes(cfg)

	Info.Println("Server running at :8080")
	http.ListenAndServe(":8080", route)
}
