package upgrade

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseName(t *testing.T) {
	tests := []struct {
		Name  string
		Ver   int
		Dir   Direction
		Error error
	}{
		{
			"1_up_foo.sql", 1, Upgrade, nil,
		},
		{
			"001_up_foo.sql", 1, Upgrade, nil,
		},
		{
			"001_dn_foo.sql", 1, Downgrade, nil,
		},
		{
			"001_down_foo.sql", 1, Downgrade, nil,
		},
		{
			"1_dn", 1, Downgrade, nil,
		},
		{
			"1_up.sql", 1, Upgrade, nil,
		},
		{
			"_up.sql", -1, Direction(-1), ErrInvalidSyntax,
		},
		{
			"1.sql", -1, Direction(-1), ErrInvalidSyntax,
		},
		{
			"1__.sql", -1, Direction(-1), ErrInvalidSyntax,
		},
		{
			"1_nope_foo.sql", -1, Direction(-1), ErrInvalidSyntax,
		},
		{
			"1_nopebutquiteabitlonger_foo.sql", -1, Direction(-1), ErrInvalidSyntax,
		},
	}
	for _, e := range tests {
		v, d, err := parseName(e.Name)
		if e.Error != nil {
			fmt.Printf("--> [%v] %s\n", err, e.Name)
			assert.Equal(t, e.Error, err, e.Name)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			fmt.Printf("--> [%v, %d] %s\n", v, d, e.Name)
			assert.Equal(t, e.Ver, v, e.Name)
			assert.Equal(t, e.Dir, d, e.Name)
		}
	}
}

// func TestValidVersions(t *testing.T) {
// }
