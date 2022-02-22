package store

import (
	"database/sql"
)

var pdb *sql.DB

type ShortURL struct {
	ShortID     string
	OriginalURL string
	ExpireTime  string
}

//Connect to postgres and redis
func Init(postgres_db *sql.DB) {
	pdb = postgres_db
}
