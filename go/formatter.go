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
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/skeema/tengo"
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
// values that know how to format a diff
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
	return nil, fmt.Errorf("Unkown formatter, only (%s) are allowed", strings.Join(existingFormatters(), ","))
}

// SQLFormatter formats a Diff in SQL format
// (ALTER, CREATE and DROP statements)
type SQLFormatter struct{}

// Format formats a diff returning a slice of string commands, each of
// which is an SQL ALTER, CREATE or DROP statement.
func (f *SQLFormatter) Format(diff *Diff) interface{} {
	return diff.Raw().String()
}

// CompactFormatter formats a diff in a compact human-readable way
type CompactFormatter struct{}

// Format returns a string with the formatted diff
func (f *CompactFormatter) Format(diff *Diff) interface{} {
	var res []string
	ods := diff.Compute()
	for _, od := range ods {
		switch od.DiffType() {
		case tengo.DiffTypeAlter:
			res = append(res, f.formatAlter(od, diff)...)
		}
	}

	return f.summarize(res)
}

func (f *CompactFormatter) formatAlter(diff tengo.ObjectDiff, context *Diff) []string {
	clauses := diff.(*TableDiff).AlterClauses()
	res := make([]string, len(clauses))
	for i, c := range clauses {
		res[i] = f.formatAlterClause(c, context)
	}
	return res
}

func (f *CompactFormatter) formatAlterClause(c tengo.TableAlterClause, context *Diff) string {
	var s string
	switch c.(type) {
	case tengo.AddColumn:
		s = f.formatAddColumn(c.(tengo.AddColumn), context)
	default:
		log.Panicf("Unexpected Table Alter Clause: %T", c)
	}
	return s
}

func (f *CompactFormatter) formatAddColumn(ac tengo.AddColumn, context *Diff) string {
	cName := ac.Column.Name
	tName := ac.Table.Name
	return fmt.Sprintf("Table %s differs: missing column %s on %s.%s", tName, cName, context.to.Name, context.to.host)
}

func (f *CompactFormatter) summarize(diffs []string) string {
	var buffer bytes.Buffer
	if count := len(diffs); count > 0 {
		buffer.WriteString(fmt.Sprintf("Differences found (%d):\n", count))
		for _, s := range diffs {
			buffer.WriteString(fmt.Sprintf("\t- %s\n", s))
		}
	} else {
		buffer.WriteString("No differences found")
	}
	return buffer.String()
}
