package sqlite3

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/bww/go-upgrade/v1"

	_ "github.com/mattn/go-sqlite3"
)

const versionTable = "schema_version"

type Driver struct {
	*sql.DB
}

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
	return d.migrateVersion(versionTable, v.Upgrade, v.Version, upgrade.Upgrade)
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

func (d *Driver) migrateVersion(t string, src io.ReadCloser, v int, dir upgrade.Direction) error {
	defer src.Close()

	sql, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if tx != nil {
			err := tx.Rollback()
			if err != nil {
				fmt.Printf("upgrade: sqlite: Could not rollback failed transaction: %v", err)
			}
		}
	}()

	_, err = tx.Exec(string(sql))
	if err != nil {
		return err
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
		return upgrade.ErrInvalidDirection
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	tx = nil

	return nil
}
