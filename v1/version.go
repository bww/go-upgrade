package upgrade

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidSyntax = errors.New("Invalid syntax")
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
//
//  <version>_<up|down>[_<optional_description>]
//
func versionsFromResourcesAtPath(p string) ([]*Version, error) {
	versions := make(map[int]*Version)

	dir, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	infos, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	for _, e := range infos {
		n := e.Name()
		if n == "" || n[0] == '.' {
			continue // skip dot files
		}

		v, d, err := parseName(n)
		if err == ErrInvalidSyntax {
			continue // not a migration file
		} else if err != nil {
			return nil, err
		}
		if v < 1 {
			return nil, fmt.Errorf("Version cannot be less than one: [%s] %v", n, v)
		}

		ver, ok := versions[v]
		if !ok {
			ver = &Version{Version: v}
			versions[v] = ver
		}

		f, err := os.Open(path.Join(p, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("Could not open resource: [%v] %v", n, err)
		}

		data, err := ioutil.ReadAll(f)
		f.Close() // must close on either success or failure
		if err != nil {
			return nil, fmt.Errorf("Could not read resource: [%v] %v", n, err)
		}

		switch d {
		case Upgrade:
			if ver.Upgrade != nil {
				return nil, fmt.Errorf("Upgrade resource redefined for version %v [%v]", ver.Version, n)
			}
			ver.Upgrade = data
		case Downgrade:
			if ver.Rollback != nil {
				return nil, fmt.Errorf("Rollback resource redefined for version %v [%v]", ver.Version, n)
			}
			ver.Rollback = data
		default:
			return nil, fmt.Errorf("Upgrade resource has invalid form: invalid resource type [%v] %v", n, v)
		}

	}

	i, out := 0, make([]*Version, len(versions))
	for _, e := range versions {
		if e.Upgrade == nil { // we only require an upgrade resource
			return nil, fmt.Errorf("Version %v is missing an upgrade resource", e.Version)
		}
		out[i] = e
		i++
	}

	sort.Sort(byVersion(out))
	return out, nil
}

func parseName(n string) (int, Direction, error) {
	n, v, err := parseVersion(n)
	if err == ErrInvalidSyntax {
		return -1, Direction(-1), ErrInvalidSyntax
	}

	if len(n) < 1 {
		return -1, Direction(-1), ErrInvalidSyntax
	}
	if !strings.ContainsRune("_-", rune(n[0])) {
		return -1, Direction(-1), ErrInvalidSyntax
	}
	if n = n[1:]; len(n) < 1 {
		return -1, Direction(-1), ErrInvalidSyntax
	}

	n, d, err := parseDirection(n)
	if err != nil {
		return -1, Direction(-1), ErrInvalidSyntax
	}

	return v, d, nil
}

func parseVersion(n string) (string, int, error) {
	x, l := 0, len(n)
	for i := 0; i < l; i++ {
		if unicode.IsDigit(rune(n[i])) {
			x++
		} else {
			break
		}
	}
	if x < 1 {
		return n, -1, ErrInvalidSyntax
	}
	v, err := strconv.Atoi(n[:x])
	if err != nil {
		return n, -1, err
	}
	return n[x:], v, nil
}

var dirnames = map[string]Direction{
	"up":   Upgrade,
	"dn":   Downgrade,
	"down": Downgrade,
}

func parseDirection(n string) (string, Direction, error) {
	x, l, m := 0, len(n), len("down") // "down" is the longest acceptable string
	for i := 0; i < l && i < m; i++ {
		if unicode.IsLetter(rune(n[i])) {
			x++
		} else {
			break
		}
	}
	d, ok := dirnames[strings.ToLower(n[:x])]
	if !ok {
		return n, Direction(-1), ErrInvalidSyntax
	}
	return n[x:], d, nil
}
