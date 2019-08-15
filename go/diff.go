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
	"reflect"
	"unsafe"

	"github.com/skeema/tengo"
)

// Diff represents a diff between schemas
type Diff struct {
	from *Schema
	to   *Schema
}

// NewDiff creates a new diff from the given schemas
func NewDiff(from, to *Schema) *Diff {
	return &Diff{
		from: from,
		to:   to,
	}
}

// Schema is a tengo.Schema enriched with Host information
type Schema struct {
	*tengo.Schema
	Host string
}

// NewSchema creates a new Schema from a tengo.Schema and
// the Host information
func NewSchema(s *tengo.Schema, host string) *Schema {
	return &Schema{
		Schema: s,
		Host:   host,
	}
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
	return res
}

// Raw returns the tengo.SchemaDiff between the receiver's
// from an to fields
func (d *Diff) Raw() *tengo.SchemaDiff {
	return tengo.NewSchemaDiff(d.from.Schema, d.to.Schema)
}

// Difference is an adapter to the tengo.ObjectDiff struct
// that exposes the alterClauses field, in order to be
// visitable by any formatter object.
//
// There's no reason why tengo.TableDiff.alterClauses is not
// exported, and a new ticket should be opened in skeema/tengo
// to export this field, and eventually remove this code, which
// is unsafe if the dependency on skeema/tengo is upgraded.
type TableDiff struct {
	*tengo.TableDiff
}

// AlterClauses returns the unexported alterClauses field of the
// adapted tengo.tableDiff so a Formatter can visit them to
// generate instructions for amending the schemas, not only in SQL
// but in the different formats supported by mydiff.
//
// TableAlterClause can be one of:
//	AddColumn
//	DropColumn
//	AddIndex
//	DropIndex
//	AddForeignKey
//	DropForeignKey
//	RenameColumn
//	ModifyColumn
//	ChangeAutoIncrement
//	ChangeCharSet
//	ChangeCreateOptions
//	ChangeComment
//	ChangeStorageEngine
func (d *TableDiff) AlterClauses() []tengo.TableAlterClause {
	val := reflect.ValueOf(d.TableDiff).Elem()
	f := val.FieldByName("alterClauses")

	// As said above, There's no reason why tengo.TableDiff.alterClauses
	// is not exported, but until that's fixed in tengo, we return an unsafe
	// reference to the unexported field.
	return *reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Interface().(*[]tengo.TableAlterClause)
}
