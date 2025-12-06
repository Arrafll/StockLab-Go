package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Arrafll/StockLab-Go/internal/config"
	_ "github.com/lib/pq"
)

// Global DB variable
var DB *sql.DB

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	// Inisialisasi variable global
	DB = db

	return db, nil
}
