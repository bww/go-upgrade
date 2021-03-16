//+build sqlite3

package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New(u string) (*Driver, error) {
	db, err := sql.Open("sqlite3", u)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return NewWithDB(db)
}
