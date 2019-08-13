// mydiff - Compute the differences between two MySQL schemas.
//
// Copyright (c) 2019 Miguel Fernández Fernández
//
// This Source Code Form is subject to the terms of MIT License:
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
	"sql":    &SQLFormatter{},
	"gh-ost": &GhostFormatter{},
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

// FormatOptions is a raw string that can modify the behaviour
// of a formatter. Different formatters will parse this
// string and validate it
type FormatOptions string

// NoFormatOptions is used to pass no options
const NoFormatOptions = FormatOptions("")

// Formatter is the interface implemented by different
// values that know how to format a diff
type Formatter interface {
	Format(diff *Diff, opts FormatOptions) interface{}
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
	return nil, fmt.Errorf("Unkown formatter, only (%s) are allowed", strings.Join(existingFormatters(), ","))
}

// SQLFormatter formats a Diff in SQL format
// (ALTER, CREATE and DROP statements)
type SQLFormatter struct{}

// Format formats a diff returning a slice of string commands, each of
// which is an SQL ALTER, CREATE or DROP statement.
func (f *SQLFormatter) Format(diff *Diff, _ FormatOptions) interface{} {
	return diff.Raw().String()
}

// GhostFormatter formats a diff in a set of gh-ost scripts
//
// See also the following gh-ost resources:
// * Cheatsheet: https://github.com/github/gh-ost/blob/master/doc/cheatsheet.md
// * Supported DDL and Limitations: https://github.com/github/gh-ost/blob/master/doc/shared-key.md
// * Command line flags: https://github.com/github/gh-ost/blob/master/doc/command-line-flags.md
//
type GhostFormatter struct{}

// Format formats a diff returning a slice of string commands, each of
// which is an SQL ALTER, CREATE or DROP statement.
func (f *GhostFormatter) Format(diff *Diff, opts FormatOptions) interface{} {
	return []string{}
}
