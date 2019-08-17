package mydiff

import (
	"github.com/go-sql-driver/mysql"
)

type ParsedDSN struct {
	*mysql.Config
}

// ParseDSN parses a DSN string into an object capable
// of returning information about its parts
func ParseDSN(DSN string) *ParsedDSN {
	c, err := mysql.ParseDSN(DSN)
	if err != nil {
		panic(err)
	}
	return &ParsedDSN{c}
}
