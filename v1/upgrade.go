package upgrade

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNoVersions = errors.New("No versions")
	ErrNoDriver   = errors.New("No driver")
	ErrNoChange   = errors.New("No change")
)

type Results struct {
	Before, After, Target int
	Applied               []int
}

func (r Results) String() string {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("Upgraded to %d (%d → %d) applied ", r.Target, r.Before, r.After))
	if len(r.Applied) > 0 {
		b.WriteString("[")
		for i, e := range r.Applied {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(strconv.FormatInt(int64(e), 10))
		}
		b.WriteString("]")
	} else {
		b.WriteString("none")
	}
	return b.String()
}

type Config struct {
	Resources string
	Driver    Driver
}

type Driver interface {
	// Obtain the current version
	Version() (int, error)
	// Execute an upgrade to the provided version
	Upgrade(Version) error
}

type Upgrader struct {
	versions []*Version
	driver   Driver
}

func New(conf Config) (*Upgrader, error) {
	v, err := versionsFromResourcesAtPath(conf.Resources)
	if err != nil {
		return nil, err
	}
	if conf.Driver == nil {
		return nil, ErrNoDriver
	}
	return &Upgrader{v, conf.Driver}, nil
}

// Upgrade to the latest version
func (u *Upgrader) Upgrade() (Results, error) {
	return u.UpgradeToVersion(-1)
}

// Upgrade to a specific version. If the version set is sparse and the requested
// version is not represented, upgrades will be performed up to the last version
// before the requested version.
//
// Specifically, if versions {1, 2, 5} are defined and version 4 is requested,
// upgrades will be performed for versions {1, 2}.
func (u *Upgrader) UpgradeToVersion(v int) (Results, error) {
	if len(u.versions) < 1 {
		return Results{}, ErrNoVersions
	}

	before, err := u.driver.Version()
	if err != nil {
		return Results{}, err
	}

	if v < 0 {
		v = u.versions[len(u.versions)-1].Version
	}
	if before >= v {
		return Results{Before: before, After: before, Target: v}, nil
	}

	after := before
	var applied []int
	for _, e := range u.versions {
		if e.Version <= after {
			continue
		} else if e.Version > v {
			break
		}

		err = u.driver.Upgrade(*e)
		if err != nil {
			return Results{Before: before, After: after, Target: v, Applied: applied}, err
		}

		after = e.Version
		applied = append(applied, e.Version)
	}

	return Results{Before: before, After: after, Target: v, Applied: applied}, nil
}
