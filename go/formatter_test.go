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
	sql1 := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id INT AUTO_INCREMENT,
			title CHAR(255) NOT NULL,
			PRIMARY KEY (id)
		)  ENGINE=INNODB;`,
	}

	sql2 := []string{
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

	s1Name, s2Name := TestCluster.LoadSchemas(t, sql1, sql2)
	from := NewServer1Schema(s1Name)
	to := NewServer2Schema(s2Name)

	diff := NewDiff(from, to)
	sqlFmt, _ := NewFormatter("sql")
	sql := sqlFmt.Format(diff)

	expected := `ALTER TABLE "tasks" MODIFY COLUMN "id" bigint(20) NOT NULL AUTO_INCREMENT, MODIFY COLUMN "title" varchar(255) NOT NULL, ADD COLUMN "owner_id" int(11) DEFAULT NULL;
CREATE TABLE "owners" (
  "id" int(11) NOT NULL AUTO_INCREMENT,
  "name" varchar(255) NOT NULL,
  PRIMARY KEY ("id")
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
`
	expected = strings.ReplaceAll(expected, "\"", "`")
	Equal(t, expected, sql)
}
