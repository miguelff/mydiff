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

	"github.com/skeema/tengo"

	. "gotest.tools/assert"
)

func TestDiff(t *testing.T) {
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
	objectDiffs := diff.Compute()
	Equal(t, 2, len(objectDiffs))

	alters := objectDiffs[0].(*TableDiff).AlterClauses()
	Equal(t, "MODIFY COLUMN `id` bigint(20) NOT NULL AUTO_INCREMENT", alters[0].Clause(tengo.StatementModifiers{}))
	Equal(t, "MODIFY COLUMN `title` varchar(255) NOT NULL", alters[1].Clause(tengo.StatementModifiers{}))
	Equal(t, "ADD COLUMN `owner_id` int(11) DEFAULT NULL", alters[2].Clause(tengo.StatementModifiers{}))

	create, _ := objectDiffs[1].(*TableDiff).Statement(tengo.StatementModifiers{})
	Assert(t, strings.HasPrefix(create, "CREATE TABLE `owners`"))
}
