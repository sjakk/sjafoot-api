package main

import (
	"log"
	"os"
	"github.com/sjakk/sjafoot/internal/data"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var testApp application

func TestMain(m *testing.M) {
	var cfg config
	cfg.port = 4001
	cfg.env = "testing"
	cfg.jwt.secret = "a-test-secret-that-is-super-secure"
	cfg.db.maxOpenConns = 25
	cfg.db.maxIdleConns = 25
	cfg.db.maxIdleTime = "15m"

	cfg.db.dsn = os.Getenv("SJAFOOT_DB_DSN")
	if cfg.db.dsn == "" {
		cfg.db.dsn = "postgres://sjafoot_user:yourpassword@localhost/sjafoot_test?sslmode=disable"
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatalf("Could not connect to test database: %s", err)
	}

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatalf("cannot create migration driver: %v", err)
	}
	migrator, err := migrate.NewWithDatabaseInstance("file://../../migrations", "postgres", migrationDriver)
	if err != nil {
		logger.Fatalf("cannot create migrator: %v", err)
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatalf("migration failed: %v", err)
	}
	logger.Println("test database migrations applied")

	testApp = application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	os.Exit(m.Run())
}
