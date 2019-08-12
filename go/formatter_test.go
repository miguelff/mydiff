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
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestNewFormatter_Existing(t *testing.T) {
	formatter, _ := NewFormatter("sQl")
	IsType(t, formatter, &SQLFormatter{})
}

func TestNewFormatter_NonExisting(t *testing.T) {
	formatter, err := NewFormatter("foo")
	Error(t, err, "fasfasdf")
	Nil(t, formatter)
}

func TestSQLFormatter_Format(t *testing.T) {
	t.Skip()
}
