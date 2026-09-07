package db

import (
  "database/sql"
  "os"

  _ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
  url := os.Getenv("DB_URL")
  return sql.Open("postgres", url)
}
