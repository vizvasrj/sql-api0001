package database

import (
	"database/sql"
	"fmt"
	"os"
	"src/conf"

	_ "github.com/lib/pq"
)

func Get_database() (*sql.DB, error) {
	pg_host := os.Getenv("POSTGRES_HOST")
	pg_db := os.Getenv("POSTGRES_DB")
	pg_user := os.Getenv("POSTGRES_USER")
	pg_pass := os.Getenv("POSTGRES_PASSWORD")
	db_url := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		pg_user, pg_pass, pg_host, pg_db,
	)
	// fmt.Printf("%s\n", db_url)
	db, err := sql.Open("postgres", db_url)
	// defer db.Close()
	if err != nil {
		return &sql.DB{}, conf.AddCallerInfo(err)
	}

	err = db.Ping()
	if err != nil {
		return &sql.DB{}, conf.AddCallerInfo(err)
	}

	return db, nil
}
