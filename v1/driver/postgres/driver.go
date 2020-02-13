package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bww/go-upgrade/v1"

	"github.com/lib/pq"
)

const versionTable = "schema_version"

var ErrInvalidDirection = errors.New("Invalid direction")

type Driver struct {
	*sql.DB
}

func New(u string) (*Driver, error) {
	db, err := sql.Open("postgres", u)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return NewWithDB(db)
}

func NewWithDB(db *sql.DB) (*Driver, error) {
	d := &Driver{db}

	err := d.createVersionTableIfNecessary(versionTable)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Driver) Version() (int, error) {
	return d.databaseVersion(versionTable)
}

func (d *Driver) Upgrade(v upgrade.Version) error {
	return d.migrateVersion(versionTable, string(v.Upgrade), v.Version, upgrade.Upgrade)
}

func (d *Driver) createVersionTableIfNecessary(t string) error {
	_, err := d.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (version int not null primary key);", t))
	if err != nil {
		return err
	}
	return nil
}

func (d *Driver) databaseVersion(t string) (int, error) {
	var version int
	err := d.QueryRow(fmt.Sprintf("SELECT version FROM %s ORDER BY version DESC LIMIT 1", t)).Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return version, nil
}

func (d *Driver) migrateVersion(t, s string, v int, dir upgrade.Direction) error {

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if tx != nil {
			err := tx.Rollback()
			if err != nil {
				fmt.Printf("upgrade: postgres: Could not rollback failed transaction: %v", err)
			}
		}
	}()

	_, err = tx.Exec(s)
	if err != nil {
		perr, ok := err.(*pq.Error)
		if ok {
			return fmt.Errorf("%s %v: %s\n\t%s", perr.Severity, perr.Code, perr.Message, perr.Detail)
		} else {
			return err
		}
	}

	switch dir {
	case upgrade.Upgrade:
		_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", t), v)
		if err != nil {
			return err
		}
	case upgrade.Downgrade:
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE version = $1", t), v)
		if err != nil {
			return err
		}
	default:
		return ErrInvalidDirection
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	tx = nil

	return nil
}

func CreateDatabase(u, n string) error {

	db, err := sql.Open("postgres", u)
	if err != nil {
		return err
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		return err
	}

	var exists int
	err = db.QueryRow(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", n)).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if exists == 1 {
		return nil // ok
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, n))
	if err != nil {
		return err
	}

	return nil
}

func DropDatabase(u, n string) error {

	db, err := sql.Open("postgres", u)
	if err != nil {
		return err
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, n))
	if err != nil {
		return err
	}

	return nil
}
