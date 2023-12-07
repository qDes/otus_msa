package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"otus_msa_docker/internal/http"
	"otus_msa_docker/internal/users"
)

/*
func health(w http.ResponseWriter, req *http.Request) {
	resp := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

	http.HandleFunc("/health", health)
	http.ListenAndServe(":8000", nil)

*/

func main() {
	// Setup database connection
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5444/postgres?sslmode=disable")
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
