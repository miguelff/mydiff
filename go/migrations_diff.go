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
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/skeema/tengo"
)

// MigrationDiff is an implementation of a tengo.ObjectDiff
// aimed at representing the difference in the migrations recorded
// in both schemas.
//
// MigrationsDiff is tested in integration in diff_test.go
// and formatter_test.go
type MigrationsDiff struct {
	Context       *Diff
	Table, Column string
	Missing1      []string
	Missing2      []string
}

// ComputeMigrationsDiff calculates a MigrationsDiff Object, which represents
// the differences between two tables containing the versions of the migrations
// that were run in two servers.
//
// This is useful while detecting inconsistencies
// in the DBs of web application development frameworks such as rails
func NewMigrationsDiff(d *Diff) (m *MigrationsDiff, err error) {
	var table, col string
	parts := strings.Split(d.MigrationsCol, ".")
	if len(parts) == 2 {
		table = parts[0]
		col = parts[1]
	}

	m = &MigrationsDiff{
		Context:  d,
		Table:    table,
		Column:   col,
		Missing1: []string{},
		Missing2: []string{},
	}

	dsn1 := *d.DSN1
	dsn1.DBName = d.From.Name
	migrations1, err := m.existingMigrations(dsn1, col, table)
	if err != nil {
		log.Warningf("Cannot retrieve migrations from %s.%s in %s/%s. Error: %s", col, table, dsn1.Addr, dsn1.DBName, err)
		return
	}

	dsn2 := *d.DSN2
	dsn2.DBName = d.To.Name
	migrations2, err := m.existingMigrations(dsn2, col, table)
	if err != nil {
		log.Warningf("Cannot retrieve migrations from %s.%s in %s/%s. Error: %s", col, table, dsn1.Addr, dsn1.DBName, err)
		return
	}

	m.Missing1 = StringSetDiff(migrations2, migrations1)
	m.Missing2 = StringSetDiff(migrations1, migrations2)
	return
}

// DiffType (see tengo.ObjectType)
func (m *MigrationsDiff) DiffType() tengo.DiffType {
	return DiffTypeMigrations
}

// ObjectKey (see tengo.ObjectType)
func (m *MigrationsDiff) ObjectKey() tengo.ObjectKey {
	return tengo.ObjectKey{
		Type: tengo.ObjectTypeTable,
		Name: m.Table,
	}
}

// Statement (see tengo.ObjectType)
func (m *MigrationsDiff) Statement(tengo.StatementModifiers) (string, error) {
	panic("Not implemented yet")
}

// IsEmpty determines whether the migrations diff is empty
func (m *MigrationsDiff) IsEmpty() bool {
	return len(m.Missing1) == 0 && len(m.Missing2) == 0
}

func (m *MigrationsDiff) existingMigrations(DSN ParsedDSN, col string, table string) ([]string, error) {
	db, err := sql.Open("mysql", DSN.FormatDSN())
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM %s ORDER BY %s", col, table, col))
	if err != nil {
		return nil, err
	}

	migrations := []string{}
	for rows.Next() {
		var migration string
		err = rows.Scan(&migration)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}
	return migrations, nil
}
