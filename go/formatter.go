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
	"strings"

	log "github.com/sirupsen/logrus"

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
	Origin interface{}
	Text   string
}

// ignoredLine represents a line that is ignored by the formatter.
// As an example when a column auto-increment varies between schemas,
// tengo represents it with two alter clauses, one with the column
// definition change, and the other one representing the own auto_increment
// change. The latter is ignored by this formatter.
var ignoredLine = line{}

// Format returns a string with the formatted diff
func (f *CompactFormatter) Format(diff *Diff) interface{} {
	var lines []line
	ods := diff.Compute()
	for _, od := range ods {
		switch od.DiffType() {
		case tengo.DiffTypeAlter:
			lines = append(lines, f.formatAlter(od, diff)...)
		case tengo.DiffTypeCreate:
			lines = append(lines, f.formatCreate(od.(*TableDiff), diff))
		case tengo.DiffTypeDrop:
			lines = append(lines, f.formatDrop(od.(*TableDiff), diff))
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
		if fa == ignoredLine {
			continue
		}
		switch fa.Origin.(type) {
		case tengo.AddForeignKey:
			s = s[:len(s)-1]
		case tengo.DropForeignKey:
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
	var l line
	switch c.(type) {
	case tengo.AddColumn:
		l = line{
			Text:   f.formatAddColumn(c.(tengo.AddColumn), context, tableName),
			Origin: c,
		}
	case tengo.DropColumn:
		l = line{
			Text:   f.formatDropColumn(c.(tengo.DropColumn), context, tableName),
			Origin: c,
		}
	case tengo.AddIndex:
		l = line{
			Text:   f.formatAddIndex(c.(tengo.AddIndex), context, tableName),
			Origin: c,
		}
	case tengo.DropIndex:
		l = line{
			Text:   f.formatDropIndex(c.(tengo.DropIndex), context, tableName),
			Origin: c,
		}
	case tengo.AddForeignKey:
		l = line{
			Text:   f.formatAddForeignKey(c.(tengo.AddForeignKey), context, tableName),
			Origin: c,
		}
	case tengo.DropForeignKey:
		l = line{
			Text:   f.formatDropForeignKey(c.(tengo.DropForeignKey), context, tableName),
			Origin: c,
		}
	case tengo.ModifyColumn:
		l = line{
			Text:   f.formatModifyColumn(c.(tengo.ModifyColumn), context, tableName),
			Origin: c,
		}
	case tengo.ChangeCharSet:
		l = line{
			Text:   f.formatChangeCharset(c.(tengo.ChangeCharSet), context, tableName),
			Origin: c,
		}
	case tengo.ChangeAutoIncrement:
		// information to render an autoincrement in change came already in a previous
		// ModifyColumn alter clause
		l = ignoredLine
	default:
		log.Errorf("Unexpected Table Alter Clause in Compact Formatter: %T. Ignoring", c)
		l = ignoredLine
	}
	return l
}

func (f *CompactFormatter) formatAddColumn(ac tengo.AddColumn, context *Diff, tableName string) string {
	return fmt.Sprintf("Table %s differs: missing column %s in %s.%s", tableName, ac.Column.Name, context.from.Name, context.from.Host)
}

func (f *CompactFormatter) formatDropColumn(dc tengo.DropColumn, context *Diff, tableName string) string {
	return fmt.Sprintf("Table %s differs: missing column %s in %s.%s", tableName, dc.Column.Name, context.to.Name, context.to.Host)
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
	return fmt.Sprintf("Table %s differs: missing %s %s(%s) in %s.%s", tableName, idxType, idxName, strings.Join(colNames, ", "), context.from.Name, context.from.Host)
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
	return fmt.Sprintf("Table %s differs: missing %s %s(%s) in %s.%s", tableName, idxType, idxName, strings.Join(colNames, ", "), context.to.Name, context.to.Host)
}

func (f *CompactFormatter) formatAddForeignKey(key tengo.AddForeignKey, context *Diff, tableName string) string {
	fkName := key.ForeignKey.Name
	colNames := make([]string, len(key.ForeignKey.Columns))
	for i, c := range key.ForeignKey.Columns {
		colNames[i] = c.Name
	}
	refName := key.ForeignKey.ReferencedTableName
	refColNames := key.ForeignKey.ReferencedColumnNames
	return fmt.Sprintf("Table %s differs: missing FOREIGN KEY %s(%s) REFERENCES %s(%s) in %s.%s", tableName, fkName, strings.Join(colNames, ", "), refName, strings.Join(refColNames, ","), context.from.Name, context.from.Host)
}

func (f *CompactFormatter) formatDropForeignKey(key tengo.DropForeignKey, context *Diff, tableName string) string {
	fkName := key.ForeignKey.Name
	colNames := make([]string, len(key.ForeignKey.Columns))
	for i, c := range key.ForeignKey.Columns {
		colNames[i] = c.Name
	}
	refName := key.ForeignKey.ReferencedTableName
	refColNames := key.ForeignKey.ReferencedColumnNames
	return fmt.Sprintf("Table %s differs: missing FOREIGN KEY %s(%s) REFERENCES %s(%s) in %s.%s", tableName, fkName, strings.Join(colNames, ", "), refName, strings.Join(refColNames, ","), context.to.Name, context.to.Host)
}

func (f *CompactFormatter) formatModifyColumn(mc tengo.ModifyColumn, context *Diff, tableName string) string {
	colName := mc.OldColumn.Name
	s1ColDef := f.colDef(mc.OldColumn)
	s2ColDef := f.colDef(mc.NewColumn)
	if s1ColDef != s2ColDef {
		return fmt.Sprintf("Table %s differs: column %s differs in column type: %s in %s.%s, %s in %s.%s", tableName, colName, s1ColDef, context.from.Name, context.from.Host, s2ColDef, context.to.Name, context.to.Host)
	}
	return fmt.Sprintf("Table %s differs: column %s AUTO_INCREMENT value differs between  %s.%s, and %s.%s", tableName, colName, context.from.Name, context.from.Host, context.to.Name, context.to.Host)
}

func (f *CompactFormatter) colDef(c *tengo.Column) string {
	// TODO: get the flavor from the context, based on information gathered from the instances.
	colDef := c.Definition(tengo.FlavorUnknown, nil)
	colDef = strings.Replace(colDef, "`"+c.Name+"` ", "", 1)
	return colDef
}

func (f *CompactFormatter) formatChangeCharset(set tengo.ChangeCharSet, context *Diff, tableName string) string {
	return fmt.Sprintf("Table %s differs: encoding changed to %s in %s.%s", tableName, set.Clause(tengo.StatementModifiers{}), context.to.Name, context.to.Host)
}

func (f *CompactFormatter) formatCreate(td *TableDiff, context *Diff) line {
	return line{
		Text:   fmt.Sprintf("Table %s is absent in %s.%s", td.To.Name, context.from.Name, context.from.Host),
		Origin: tengo.DiffTypeCreate,
	}
}

func (f *CompactFormatter) formatDrop(td *TableDiff, context *Diff) line {
	return line{
		Text:   fmt.Sprintf("Table %s is absent in %s.%s", td.From.Name, context.to.Name, context.to.Host),
		Origin: tengo.DiffTypeCreate,
	}
}
