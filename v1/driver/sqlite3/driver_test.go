//+build sqlite

package sqlite3

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/bww/go-upgrade/v1"

	"github.com/stretchr/testify/assert"
)

const (
	testDB = "upgrade_sqlite3.db"
)

func TestUpgrade(t *testing.T) {
	dbpath := path.Join(os.TempDir(), testDB)
	fmt.Println("-->", dbpath)
	os.Remove(dbpath)

	t.Run("a", func(t *testing.T) {
		var r upgrade.Results
		var err error
		var n int

		d, err := New(dbpath)
		assert.Nil(t, err, fmt.Sprintf("%v", err))
		defer func() {
			d.Close()
			os.Remove(dbpath)
		}()

		u, err := upgrade.New(upgrade.Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "generic/001"), Driver: d})
		assert.Nil(t, err, fmt.Sprintf("%v", err))

		r, err = u.Upgrade()
		if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			assert.Equal(t, 2, r.After)
		}

		n, err = d.Version()
		if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			assert.Equal(t, 2, n)
		}

		u, err = upgrade.New(upgrade.Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "generic/002"), Driver: d})
		assert.Nil(t, err, fmt.Sprintf("%v", err))

		r, err = u.Upgrade()
		assert.Nil(t, err, fmt.Sprintf("%v", err))
		assert.Equal(t, 3, r.After)

		n, err = d.Version()
		assert.Nil(t, err, fmt.Sprintf("%v", err))
		assert.Equal(t, 3, n)

	})

}
