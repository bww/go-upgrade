package upgrade

import (
	"fmt"
)

/**
 * Upgrader configuration
 */
type Config struct {
	Resources string
	Driver    Driver
}

/**
 * Implemented by upgrade drivers
 */
type Driver interface {
	// Obtain the current version
	Version() (int, error)
	// Execute an upgrade to the provided version
	Upgrade(Version) error
}

/**
 * An upgrader
 */
type Upgrader struct {
	versions []*Version
	driver   Driver
}

/**
 * Create an upgrader
 */
func New(conf Config) (*Upgrader, error) {
	v, err := versionsFromResourcesAtPath(conf.Resources)
	if err != nil {
		return nil, err
	}
	if conf.Driver == nil {
		return nil, fmt.Errorf("Driver is nil")
	}
	return &Upgrader{v, conf.Driver}, nil
}

/**
 * Upgrade to the latest version
 */
func (u *Upgrader) Upgrade() (int, error) {
	return u.UpgradeToVersion(-1)
}

/**
 * Upgrade to a specific version. If the version set is sparse and the requested
 * version is not represented, upgrades will be performed up to the last version
 * before the requested version.
 *
 * Specifically, if versions {1, 2, 5} are defined and version 4 is requested,
 * upgrades will be performed for versions {1, 2}.
 */
func (u *Upgrader) UpgradeToVersion(v int) (int, error) {
	if len(u.versions) < 1 {
		return -1, fmt.Errorf("No versions")
	}

	c, err := u.driver.Version()
	if err != nil {
		return -1, err
	}

	if v < 0 {
		v = u.versions[len(u.versions)-1].Version
	}
	if c >= v {
		return c, nil
	}

	for _, e := range u.versions {
		if e.Version <= c {
			continue
		} else if e.Version > v {
			break
		}

		err = u.driver.Upgrade(*e)
		if err != nil {
			return c, err
		}

		c = e.Version
	}

	return c, nil
}
