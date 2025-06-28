package db

import (
	"fmt"
	"log"
	"time"
	"trackly-backend/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type scopeFn = func(db *gorm.DB) *gorm.DB

type CommonScopeOption = scopeFn

func InitDB(cfg *config.DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName)
	println("my db" + dsn)
	println("my db" + dsn)
	retryConnectCount := 5
	var db *gorm.DB
	var conErr error
	for i := 0; i < retryConnectCount; i++ {
		database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			conErr = err
			time.Sleep(5 * time.Second)
		}
		db = database
		conErr = nil

	}
	if conErr != nil {
		return nil, conErr
	}

	migraionDns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

	m, err := migrate.New(
		"file://migrations",
		migraionDns,
	)

	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	// Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
	return db, nil
}
