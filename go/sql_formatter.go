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

// SQLFormatter formats a Diff in SQL format
// (ALTER, CREATE and DROP statements)
type SQLFormatter struct{}

// Format formats a diff returning a slice of string commands, each of
// which is an SQL ALTER, CREATE or DROP statement.
func (f *SQLFormatter) Format(diff *Diff) interface{} {
	return diff.Raw().String()
}
