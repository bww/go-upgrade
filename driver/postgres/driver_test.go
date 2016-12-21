package postgres

import (
  "fmt"
  "testing"
  "../.."
  "github.com/stretchr/testify/assert"
)

const testDatabase = "go_upgrade_tests"

func TestUpgrade(t *testing.T) {
  
  err := createDatabase("postgres://postgres@localhost:5432/template1?sslmode=disable", testDatabase)
  if !assert.Nil(t, err, fmt.Sprintf("%v", err)) { return }
  
  d, err := New(fmt.Sprintf("postgres://postgres@localhost:5432/%s?sslmode=disable", testDatabase))
  if !assert.Nil(t, err, fmt.Sprintf("%v", err)) { return }
  
  u, err := upgrade.New(upgrade.Config{Resources:"./test/postgres/001", Driver:d})
  if !assert.Nil(t, err, fmt.Sprintf("%v", err)) { return }
  
  n, err := u.Upgrade()
  if !assert.Nil(t, err, fmt.Sprintf("%v", err)) { return }
  assert.Equal(t, 2, n)
  
  n, err = u.Upgrade()
  if !assert.Nil(t, err, fmt.Sprintf("%v", err)) { return }
  assert.Equal(t, 4, n)
  
  n, err = u.Upgrade()
  if !assert.Nil(t, err, fmt.Sprintf("%v", err)) { return }
  assert.Equal(t, 4, n)
  
}
