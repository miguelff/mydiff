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

// line corresponds to a formatter alter clause, that not necessary corresponds
// 1:1 to a schema difference. Some of these formatted alter clauses are omitted.
//
// tengo.Diff exposes the differences as alter clauses, some differences
// like the Added Foreign keys are represented by two alter clauses:
// first, an ADD KEY, and then ADD CONSTRAINT over that key.
// `combine([]line) []string` combines M formatted alters into N <=M
// strings each of which will be a difference outputted by the formatter.
type line struct {
	Type interface{}
	Text string
}

// Format returns a string with the formatted diff
func (f *CompactFormatter) Format(diff *Diff) interface{} {
	var lines []line
	ods := diff.Compute()
	for _, od := range ods {
		switch od.DiffType() {
		case tengo.DiffTypeAlter:
			lines = append(lines, f.formatAlter(od, diff)...)
		}
	}
	return f.summarize(lines)
}

// combine combines several line together into a list
// of strings each of which is a line outputted by the formatter.
//
// AddForeignKey comes in two different alter clauses:
// ADD KEY k followed by an ADD CONSTRAINT on k.
// We only care about the last one, so we pop the previous line.
func (f *CompactFormatter) combine(lines []line) (s []string) {
	for _, fa := range lines {
		switch fa.Type.(type) {
		case tengo.AddForeignKey:
			s = s[:len(s)-1]
		}
		s = append(s, fa.Text)
	}
	return s
}

func (f *CompactFormatter) summarize(fas []line) string {
	lines := f.combine(fas)
	var buffer bytes.Buffer
	if count := len(lines); count > 0 {
		buffer.WriteString(fmt.Sprintf("Differences found (%d):\n", count))
		for _, s := range lines {
			buffer.WriteString(fmt.Sprintf("\t- %s\n", s))
		}
	} else {
		buffer.WriteString("No differences found")
	}
	return buffer.String()
}

func (f *CompactFormatter) formatAlter(diff tengo.ObjectDiff, context *Diff) []line {
	tableDiff := diff.(*TableDiff)
	tableName := tableDiff.From.Name

	clauses := tableDiff.AlterClauses()
	lines := make([]line, len(clauses))
	for i, c := range clauses {
		lines[i] = f.formatAlterClause(c, context, tableName)
	}
	return lines
}

func (f *CompactFormatter) formatAlterClause(c tengo.TableAlterClause, context *Diff, tableName string) line {
	var cd line
	switch c.(type) {
	case tengo.AddColumn:
		cd = line{
			Text: f.formatAddColumn(c.(tengo.AddColumn), context, tableName),
			Type: &tengo.AddColumn{},
		}
	case tengo.DropColumn:
		cd = line{
			Text: f.formatDropColumn(c.(tengo.DropColumn), context, tableName),
			Type: tengo.DropColumn{},
		}
	case tengo.AddIndex:
		cd = line{
			Text: f.formatAddIndex(c.(tengo.AddIndex), context, tableName),
			Type: tengo.AddIndex{},
		}
	case tengo.DropIndex:
		cd = line{
			Text: f.formatDropIndex(c.(tengo.DropIndex), context, tableName),
			Type: tengo.DropIndex{},
		}
	case tengo.AddForeignKey:
		cd = line{
			Text: f.formatAddForeignKey(c.(tengo.AddForeignKey), context, tableName),
			Type: tengo.AddForeignKey{},
		}
	default:
		log.Panicf("Unexpected Table Alter Clause: %T", c)
	}
	return cd
}

func (f *CompactFormatter) formatAddColumn(ac tengo.AddColumn, context *Diff, tableName string) string {
	return fmt.Sprintf("Table %s differs: missing column %s on %s.%s", tableName, ac.Column.Name, context.from.Name, context.from.host)
}

func (f *CompactFormatter) formatDropColumn(dc tengo.DropColumn, context *Diff, tableName string) string {
	return fmt.Sprintf("Table %s differs: missing column %s on %s.%s", tableName, dc.Column.Name, context.to.Name, context.to.host)
}

func (f *CompactFormatter) formatAddIndex(idx tengo.AddIndex, context *Diff, tableName string) string {
	var idxType string
	if idx.Index.Unique {
		idxType = "UNIQUE KEY"
	} else {
		idxType = "KEY"
	}
	idxName := idx.Index.Name
	colNames := make([]string, len(idx.Index.Columns))
	for i, c := range idx.Index.Columns {
		colNames[i] = c.Name
	}
	return fmt.Sprintf("Table %s differs: missing %s %s(%s) on %s.%s", tableName, idxType, idxName, strings.Join(colNames, ", "), context.from.Name, context.from.host)
}

func (f *CompactFormatter) formatDropIndex(idx tengo.DropIndex, context *Diff, tableName string) string {
	var idxType string
	if idx.Index.Unique {
		idxType = "UNIQUE KEY"
	} else {
		idxType = "KEY"
	}
	colNames := make([]string, len(idx.Index.Columns))
	for i, c := range idx.Index.Columns {
		colNames[i] = c.Name
	}
	idxName := idx.Index.Name
	return fmt.Sprintf("Table %s differs: missing %s %s(%s) on %s.%s", tableName, idxType, idxName, strings.Join(colNames, ", "), context.to.Name, context.to.host)
}

func (f *CompactFormatter) formatAddForeignKey(key tengo.AddForeignKey, context *Diff, tableName string) string {
	fkName := key.ForeignKey.Name
	colNames := make([]string, len(key.ForeignKey.Columns))
	for i, c := range key.ForeignKey.Columns {
		colNames[i] = c.Name
	}
	refName := key.ForeignKey.ReferencedTableName
	refColNames := key.ForeignKey.ReferencedColumnNames
	return fmt.Sprintf("Table %s differs: missing FOREIGN KEY %s(%s) REFERENCES %s(%s) on %s.%s", tableName, fkName, strings.Join(colNames, ", "), refName, strings.Join(refColNames, ","), context.from.Name, context.from.host)
}
