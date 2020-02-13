package upgrade

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const (
	upgrade  = "up"
	rollback = "down"
)

type Direction int

const (
	Upgrade Direction = iota
	Downgrade
)

type Version struct {
	Version  int
	Upgrade  []byte
	Rollback []byte
}

type byVersion []*Version

func (s byVersion) Len() int {
	return len(s)
}

func (s byVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byVersion) Less(i, j int) bool {
	return s[i].Version < s[j].Version
}

// Load upgrade resources from a directory. An upgrade resource has
// a filename of the following form:
//  <version>_<up|down>[_<optional_description>]
func versionsFromResourcesAtPath(p string) ([]*Version, error) {
	versions := make(map[int]*Version)

	dir, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	infos, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	for _, e := range infos {
		n := e.Name()
		if len(n) > 0 {
			if n[0] == '.' {
				continue // skip dot files
			}
			if !unicode.IsDigit(rune(n[0])) {
				continue // skip files that don't begin with a digit
			}
		}

		x := strings.IndexAny(n, "_-")
		if x < 0 {
			return nil, fmt.Errorf("Upgrade resource has invalid form: [%v] (version)", e.Name())
		}

		v, err := strconv.Atoi(n[:x])
		if err != nil {
			return nil, fmt.Errorf("Upgrade resource has invalid form: invalid version number: [%v] %v", e.Name(), err)
		}
		if v < 1 {
			return nil, fmt.Errorf("Version cannot be less than one: [%v] %v", e.Name(), v)
		}

		n = n[x+1:]

		version, ok := versions[v]
		if !ok {
			version = &Version{Version: v}
			versions[v] = version
		}

		f, err := os.Open(path.Join(p, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("Could not open resource: [%v] %v", e.Name(), err)
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("Could not read resource: [%v] %v", e.Name(), err)
		}

		switch {
		case strings.HasPrefix(n, upgrade):
			if version.Upgrade != nil {
				return nil, fmt.Errorf("Upgrade resource redefined for version %v: [%v]", version.Version, e.Name())
			}
			version.Upgrade = data
		case strings.HasPrefix(n, rollback):
			if version.Rollback != nil {
				return nil, fmt.Errorf("Rollback resource redefined for version %v: [%v]", version.Version, e.Name())
			}
			version.Rollback = data
		default:
			return nil, fmt.Errorf("Upgrade resource has invalid form: invalid resource type [%v] %v", e.Name(), v)
		}

	}

	i, out := 0, make([]*Version, len(versions))
	for _, e := range versions {
		if e.Upgrade == nil {
			return nil, fmt.Errorf("Version %v is missing an upgrade resource", e.Version)
		}
		out[i] = e
		i++
	}

	sort.Sort(byVersion(out))
	return out, nil
}
