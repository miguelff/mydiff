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
	log "github.com/sirupsen/logrus"

	"github.com/skeema/tengo"
)

// A integer denoting our custom diff type for schema migrations
// We add 10^6 to iota, to avoid collisions with tengo diff types.
const DiffTypeMigrations = 10 ^ 6 + iota

// Diff encapsulates the data necessary to compute a diff between two schemas
// in servers denoted by DSN1, and DSN2
type Diff struct {
	DSN1, DSN2        *ParsedDSN
	From, To          *tengo.Schema
	IncludeMigrations bool
	MigrationsCol     string
}

// NewDiff creates a new Diff
func NewDiff(DSN1, DSN2 string, from, to *tengo.Schema, includeMigrations bool, migrationsCol string) *Diff {
	return &Diff{
		DSN1:              ParseDSN(DSN1),
		DSN2:              ParseDSN(DSN2),
		From:              from,
		To:                to,
		IncludeMigrations: includeMigrations,
		MigrationsCol:     migrationsCol,
	}
}

// Raw returns the tengo.SchemaDiff between the receiver's
// From an To fields
func (d *Diff) Raw() *tengo.SchemaDiff {
	return tengo.NewSchemaDiff(d.From, d.To)
}

// Compute computes the difference between the two schemas
// returning an Difference object.
func (d *Diff) Compute() []tengo.ObjectDiff {
	objectDiffs := d.Raw().ObjectDiffs()

	var res []tengo.ObjectDiff = make([]tengo.ObjectDiff, len(objectDiffs))
	for i, od := range objectDiffs {
		switch od.(type) {
		case *tengo.TableDiff:
			res[i] = &TableDiff{od.(*tengo.TableDiff)}
		default:
			res[i] = od
		}
	}

	if d.IncludeMigrations {
		migrationsDiff, err := NewMigrationsDiff(d)
		if err != nil {
			log.Warningf("Error while computing the migrations diff: %s", err)
		} else {
			if !migrationsDiff.IsEmpty() {
				res = append(res, migrationsDiff)
			}
		}
	}

	return res
}
