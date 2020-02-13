package postgres

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/bww/go-upgrade/v1"

	"github.com/stretchr/testify/assert"
)

const (
	testDatabase = "go_upgrade_tests"
	testURL      = "postgres://postgres@localhost:5432/template1?sslmode=disable"
)

func TestUpgrade(t *testing.T) {

	t.Run("a", func(t *testing.T) {
		var r upgrade.Results
		var err error
		var n int

		err = CreateDatabase(testURL, testDatabase)
		assert.Nil(t, err, fmt.Sprintf("%v", err))

		d, err := New(fmt.Sprintf("postgres://postgres@localhost:5432/%s?sslmode=disable", testDatabase))
		assert.Nil(t, err, fmt.Sprintf("%v", err))
		defer d.Close()

		u, err := upgrade.New(upgrade.Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "postgres/001"), Driver: d})
		assert.Nil(t, err, fmt.Sprintf("%v", err))

		r, err = u.Upgrade()
		if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			assert.Equal(t, 2, r.After)
		}

		n, err = d.Version()
		if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			assert.Equal(t, 2, n)
		}

		u, err = upgrade.New(upgrade.Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "postgres/002"), Driver: d})
		assert.Nil(t, err, fmt.Sprintf("%v", err))

		r, err = u.Upgrade()
		assert.Nil(t, err, fmt.Sprintf("%v", err))
		assert.Equal(t, 3, r.After)

		n, err = d.Version()
		assert.Nil(t, err, fmt.Sprintf("%v", err))
		assert.Equal(t, 3, n)

	})

	err := DropDatabase(testURL, testDatabase)
	assert.Nil(t, err, fmt.Sprintf("%v", err))
}
