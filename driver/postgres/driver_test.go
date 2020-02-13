package postgres

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/bww/go-upgrade"

	"github.com/stretchr/testify/assert"
)

const (
	testDatabase = "go_upgrade_tests"
	testURL      = "postgres://postgres@localhost:5432/template1?sslmode=disable"
)

func TestUpgrade(t *testing.T) {

	t.Run("a", func(t *testing.T) {
		var err error
		var n int

		err = CreateDatabase(testURL, testDatabase)
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}

		d, err := New(fmt.Sprintf("postgres://postgres@localhost:5432/%s?sslmode=disable", testDatabase))
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}
		defer d.Close()

		u, err := upgrade.New(upgrade.Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "postgres/001"), Driver: d})
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}

		n, err = u.Upgrade()
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}
		assert.Equal(t, 2, n)

		n, err = d.Version()
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}
		assert.Equal(t, 2, n)

		u, err = upgrade.New(upgrade.Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "postgres/002"), Driver: d})
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}

		n, err = u.Upgrade()
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}
		assert.Equal(t, 3, n)

		n, err = d.Version()
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}
		assert.Equal(t, 3, n)

	})

	err := DropDatabase(testURL, testDatabase)
	if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
		return
	}

}
