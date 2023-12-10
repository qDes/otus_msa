package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
	"otus_msa_docker/internal/http"
	"otus_msa_docker/internal/users"
)

func main() {

	dbHost := os.Getenv("postgresqlHost")
	dbPort := os.Getenv("postgresqlPort")
	dbUser := os.Getenv("postgresqlUsername")
	dbPassword := os.Getenv("postgresqlPassword")
	dbName := os.Getenv("postgresqlDatabase")

	// Construct the connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Setup database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Ensure the database is closed when the function ends
	defer db.Close()

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // file source URL
		"postgres",          // database name
		driver,
	)

	if err != nil {
		panic(err)
	}

	err = m.Up() // Apply all up migrations
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	store, _ := users.NewPostgresStore(db)

	s, err := http.NewServer(store)
	if err != nil {
		panic(fmt.Errorf("failed to start server: %w", err))
	}

	http.RunServer(s)
}
