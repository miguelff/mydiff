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
	"strings"
	"testing"

	"gotest.tools/assert"

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

// TestSQLFormatter_Format won't have many additional tests, as
// SQL output is provided by skeema/tengo, which is also properly
// tested.
func TestSQLFormatter_Format(t *testing.T) {
	schema1 := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id INT AUTO_INCREMENT,
			title CHAR(255) NOT NULL,
			PRIMARY KEY (id)
		)  ENGINE=INNODB;`,
	}

	schema2 := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id BIGINT AUTO_INCREMENT,
			title VARCHAR(255) NOT NULL,
			owner_id INT,
			PRIMARY KEY (id)
		)  ENGINE=INNODB;`,
		`CREATE TABLE IF NOT EXISTS owners (
			id INT AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			PRIMARY KEY (id)
		)  ENGINE=INNODB;`,
	}

	expected := `ALTER TABLE "tasks" MODIFY COLUMN "id" bigint(20) NOT NULL AUTO_INCREMENT, MODIFY COLUMN "title" varchar(255) NOT NULL, ADD COLUMN "owner_id" int(11) DEFAULT NULL;
CREATE TABLE "owners" (
  "id" int(11) NOT NULL AUTO_INCREMENT,
  "name" varchar(255) NOT NULL,
  PRIMARY KEY ("id")
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
`
	expected = strings.ReplaceAll(expected, "\"", "`")
	sqlFmt, _ := NewFormatter("sql")
	sql := RunDiff(t, schema1, schema2, sqlFmt)

	Equal(t, expected, sql)
}

func TestCompactFormatter_Format(t *testing.T) {
	tests := map[string]struct {
		schema1  []string
		schema2  []string
		expected string
	}{
		//	AddColumn
		"Add Column": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,
					owner_id INT,
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			expected: "Differences found:" +
				"- Table tasks differs: missing column owner_id on ",
		},
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
		//  CreateTable
		//  DropTable
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ghFmt, _ := NewFormatter("gh-ost")
			result := RunDiff(t, test.schema1, test.schema2, ghFmt)
			ghostScripts := result.([]string)
			assert.Equal(t, len(test.expected), len(ghostScripts))
			for i, r := range ghostScripts {
				assert.Equal(t, test.expected[i], r)
			}
		})
	}
}
