package postgres

import (
  "fmt"
  "testing"
  "../.."
  "github.com/stretchr/testify/assert"
)

func TestUpgrade(t *testing.T) {
  
  d, err := New("postgres://postgres@localhost:5432/go_upgrade_tests?sslmode=disable")
  if !assert.Nil(t, err, fmt.Sprintf("%v", err)) { return }
  
  u, err := upgrade.New(upgrade.Config{Resources:"./test/versions/001", Driver:d})
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
