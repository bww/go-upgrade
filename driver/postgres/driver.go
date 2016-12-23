package postgres

import (
  "fmt"
  "database/sql"
)

import (
  "../.."
)

import (
  "github.com/lib/pq"
)

const versionTable = "schema_version"

/**
 * A postgres driver
 */
type Driver struct {
  *sql.DB
}

/**
 * Create a new postgres driver for the provided connection URL
 */
func New(u string) (*Driver, error) {
  
  db, err := sql.Open("postgres", u)
  if err != nil {
    return nil, err
  }
  
  err = db.Ping()
  if err != nil {
    return nil, err
  }
  
  d := &Driver{db}
  
  err = d.createVersionTableIfNecessary(versionTable)
  if err != nil {
    return nil, err
  }
  
  return d, nil
}

/**
 * Obtain the current version of the database
 */
func (d *Driver) Version() (int, error) {
  return d.databaseVersion(versionTable)
}

/**
 * Execute an upgrade to the provided version
 */
func (d *Driver) Upgrade(v upgrade.Version) error {
  fmt.Println("postgres: -> version", v.Version)
  return d.execVersionScript(versionTable, string(v.Upgrade), v.Version)
}

/**
 * Create the version table, if we haven't already
 */
func (d *Driver) createVersionTableIfNecessary(t string) error {
  _, err := d.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (version int not null primary key);", t));
  if err != nil {
    return err
  }
  return nil
}

/**
 * Obtain the current database version
 */
func (d *Driver) databaseVersion(t string) (int, error) {
  var version int
  err := d.QueryRow(fmt.Sprintf("SELECT version FROM %s ORDER BY version DESC LIMIT 1", t)).Scan(&version)
  if err == sql.ErrNoRows {
    return 0, nil
  }else if err != nil {
    return 0, err
  }
  return version, nil
}

/**
 * Execute an upgrade
 */
func (d *Driver) execVersionScript(t, s string, v int) error {
  
  tx, err := d.Begin()
  if err != nil {
    return err
  }
  
  var success bool
  defer func() {
    if !success {
      err := tx.Rollback()
      if err != nil {
        fmt.Printf("postgres: Could not rollback failed transaction: %v", err)
      }
    }
  }()
  
  _, err = tx.Exec(s)
  if err != nil {
    perr, ok := err.(*pq.Error)
    if ok {
      return fmt.Errorf("%s %v: %s", perr.Severity, perr.Code, perr.Message)
    }else{
      return err
    }
  }
  
  _, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", t), v)
  if err != nil {
    return err
  }
  
  err = tx.Commit()
  if err != nil {
    return err
  }
  
  success = true
  return nil
}

/**
 * Create a database, for debugging
 */
func createDatabase(u, n string) error {
  
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
