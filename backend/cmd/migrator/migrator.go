package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		// credentials
		host, port, user, password, dbname, sslmode string
		// migrations
		migrationsPath, migrationsTable string
		// action
		up, down, drop bool
		force          int
	)

	flag.StringVar(&host, "host", "localhost", "Postgres host")
	flag.StringVar(&port, "port", "5432", "Postgres port")
	flag.StringVar(&user, "user", "", "Postgres user")
	flag.StringVar(&password, "password", "", "Postgres password")
	flag.StringVar(&dbname, "dbname", "", "Postgres database")
	flag.StringVar(&sslmode, "sslmode", "disable", "Postgres SSL mode")
	flag.StringVar(&migrationsPath, "migrations-path", "./migrations", "Path to migration files")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Migrations version table name")

	// Флаги действий
	flag.BoolVar(&up, "up", false, "Perform up migration")
	flag.BoolVar(&down, "down", false, "Perform rollback the last migration step (safe)")
	flag.BoolVar(&drop, "drop", false, "WARNING: Rerform rollback ALL migrations (drop database to zero)")
	flag.IntVar(&force, "force", 0, "Force migration to a specific version")

	flag.Parse()

	if user == "" || dbname == "" {
		log.Fatal("user and dbname are required")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&x-migrations-table=%s",
		user, password, host, port, dbname, sslmode, migrationsTable,
	)

	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// perform action
	switch {
	case up:
		err = m.Up()
	case down:
		// rolling back only one step (the last migration)
		err = m.Steps(-1)
	case drop:
		// drop database
		err = m.Down()
	case force != 0:
		err = m.Force(force)
	default:
		log.Fatal("either -up, -down, -drop, or -force flag must be provided")
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return
		}
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations applied successfully")
}
