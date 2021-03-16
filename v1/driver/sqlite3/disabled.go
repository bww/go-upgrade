//+build !sqlite3

package sqlite3

func New(u string) (*Driver, error) {
	panic("You must set the 'sqlite' build tag to use the 'sqlite3' package")
}
