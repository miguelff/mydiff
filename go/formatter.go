// mydiff - Compute the differences between two MySQL schemas.
//
// Copyright (c) 2019 Miguel Fernández Fernández
//
// This Source Code Form is subject To the terms of MIT License:
// A short and simple permissive license with conditions only
// requiring preservation of copyright and license notices.
// Licensed works, modifications, and larger works may be
// distributed under different terms and without source code.
//
// You can obtain a copy of the license here:
// https://opensource.org/licenses/MIT

package mydiff

import (
	"fmt"
	"strings"
)

// formatters contains a private map with the available
// formatters
var formatters map[string]Formatter = map[string]Formatter{
	"sql":     &SQLFormatter{},
	"compact": &CompactFormatter{},
}

// existingFormatters returns a slice of the existing formatters
// in the system
func existingFormatters() []string {
	keys := []string{}
	for k := range formatters {
		keys = append(keys, k)
	}
	return keys
}

// Formatter is the interface implemented by different
// values that know how To format a diff
type Formatter interface {
	Format(diff *Diff) interface{}
}

// NewFormatter creates a new value of a specific formatter based
// on the given difftype.
// Allowed difftypes are:
// - sql: which returns a SQLFormatter
// If the difftype is unknown, then an error is returned.
func NewFormatter(diffType string) (Formatter, error) {
	if formatter, ok := formatters[strings.ToLower(diffType)]; ok {
		return formatter, nil
	}
	return nil, fmt.Errorf("Unkown formatter '%s', only (%s) are allowed", diffType, strings.Join(existingFormatters(), ","))
}
