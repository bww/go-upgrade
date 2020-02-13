package upgrade

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDriver struct {
	version int
}

func (d *testDriver) Version() (int, error) {
	return d.version, nil
}

func (d *testDriver) Upgrade(v Version) error {
	d.version = v.Version
	return nil
}

func TestValidVersions(t *testing.T) {
	u, err := New(Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "versions/001"), Driver: &testDriver{0}})
	assert.Nil(t, err, fmt.Sprintf("%v", err))
	assert.Len(t, u.versions, 3)

	assert.Equal(t, 1, u.versions[0].Version)
	assert.Equal(t, []byte("1. Up"), u.versions[0].Upgrade)
	assert.Equal(t, []byte("1. Down"), u.versions[0].Rollback)

	assert.Equal(t, 2, u.versions[1].Version)
	assert.Equal(t, []byte("2. Up"), u.versions[1].Upgrade)
	assert.Equal(t, []byte("2. Down"), u.versions[1].Rollback)

	assert.Equal(t, 4, u.versions[2].Version) // this is weird, but valid; versions can be sparse
	assert.Equal(t, []byte("4. Up"), u.versions[2].Upgrade)
	assert.Equal(t, []byte("4. Down"), u.versions[2].Rollback)
}

func TestMalformedVersions(t *testing.T) {
	_, err := New(Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "versions/002")})
	if !assert.NotNil(t, err, fmt.Sprintf("%v", err)) {
		return
	}
}

func TestUpgrade(t *testing.T) {
	u, err := New(Config{Resources: path.Join(os.Getenv("GO_UPGRADE_TEST_RESOURCES"), "versions/001"), Driver: &testDriver{0}})
	assert.Nil(t, err, fmt.Sprint(err))

	r, err := u.UpgradeToVersion(2)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		fmt.Println("-->", r)
		assert.Equal(t, 2, r.After)
	}

	r, err = u.Upgrade()
	if assert.Nil(t, err, fmt.Sprint(err)) {
		fmt.Println("-->", r)
		assert.Equal(t, 4, r.After)
	}

	r, err = u.Upgrade()
	if assert.Nil(t, err, fmt.Sprint(err)) {
		fmt.Println("-->", r)
		assert.Equal(t, 4, r.After)
	}
}
