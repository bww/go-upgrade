package driver

import (
	"github.com/bww/go-upgrade/v1"
)

type Driver interface {
	Version() (int, error)
	Upgrade(v upgrade.Version) error
}
