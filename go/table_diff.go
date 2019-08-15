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
	"reflect"
	"unsafe"

	"github.com/skeema/tengo"
)

// TableDiff is an adapter of the tengo.TableDiff that exposes the
// alterClauses field, in order to be visitable by any formatter object.
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
