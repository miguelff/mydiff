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
		expected []string
	}{
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
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: missing column owner_id in schema1_\\d+.127.0.0.1:33060",
			},
		},
		"Drop Column": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,
					owner_id INT,
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: missing column owner_id in schema2_\\d+.127.0.0.1:33062",
			},
		},
		//	Add Index
		"Add Index": {
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
					PRIMARY KEY (id),
					KEY title_index (title)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: missing KEY title_index\\(title\\) in schema1_\\d+.127.0.0.1:33060",
			},
		},
		//	Add Unique Index
		"Add Unique Index": {
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
					PRIMARY KEY (id),
					UNIQUE KEY title_index (title)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: missing UNIQUE KEY title_index\\(title\\) in schema1_\\d+.127.0.0.1:33060",
			},
		},
		"Drop Index": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,
					PRIMARY KEY (id),
					KEY title_index (title)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,	
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: missing KEY title_index\\(title\\) in schema2_\\d+.127.0.0.1:33062",
			},
		},
		//	Drop Unique Index
		"Drop Unique Index": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,
					PRIMARY KEY (id),
					UNIQUE KEY title_index (title)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,	
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: missing UNIQUE KEY title_index\\(title\\) in schema2_\\d+.127.0.0.1:33062",
			},
		},
		//	Change Index
		"Change Index": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,
					PRIMARY KEY (id),
					KEY title_index (title)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					title VARCHAR(255) NOT NULL,	
					PRIMARY KEY (id),
					UNIQUE KEY title_index (title)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(2\\)",
				"Table tasks differs: missing KEY title_index\\(title\\) in schema2_\\d+.127.0.0.1:33062",
				"Table tasks differs: missing UNIQUE KEY title_index\\(title\\) in schema1_\\d+.127.0.0.1:33060",
			},
		},
		"Add Foreign Key": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					parent_id BIGINT NOT NULL,
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					parent_id BIGINT NOT NULL,
					PRIMARY KEY (id),
					FOREIGN KEY tasks_ibfk_1(parent_id) REFERENCES tasks(id) ON UPDATE CASCADE
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: missing FOREIGN KEY tasks_ibfk_1\\(parent_id\\) REFERENCES tasks\\(id\\) in schema1_\\d+.127.0.0.1:33060",
			},
		},
		//	Drop Foreign Key
		"Drop Foreign Key": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					parent_id BIGINT NOT NULL,
					PRIMARY KEY (id),
					FOREIGN KEY tasks_ibfk_1(parent_id) REFERENCES tasks(id) ON UPDATE CASCADE
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					parent_id BIGINT NOT NULL,
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: missing FOREIGN KEY tasks_ibfk_1\\(parent_id\\) REFERENCES tasks\\(id\\) in schema2_\\d+.127.0.0.1:33062",
			},
		},
		"Rename column": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					parent_id BIGINT NOT NULL,
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,
					parent_id INT NULL DEFAULT 0,
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: column parent_id differs in column type: bigint\\(20\\) NOT NULL in schema1_\\d+.127.0.0.1:33060, int\\(11\\) DEFAULT '0' in schema2_\\d+.127.0.0.1:33062",
			},
		},
		"Change Auto Increment": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT NOT NULL,	
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,	
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: column id differs in column type: bigint\\(20\\) NOT NULL in schema1_\\d+.127.0.0.1:33060, bigint\\(20\\) NOT NULL AUTO_INCREMENT in schema2_\\d+.127.0.0.1:33062",
			},
		},
		"Change Charset": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,	
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,	
					PRIMARY KEY (id)
				)  ENGINE=INNODB CHARACTER SET utf8mb4;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks differs: encoding changed To DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci in schema2_\\d+.127.0.0.1:33062",
			},
		},
		"Drop Table": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,	
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			schema2: []string{},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks is absent in schema2_\\d+.127.0.0.1:33062",
			},
		},
		"Create Table": {
			schema1: []string{},
			schema2: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id BIGINT AUTO_INCREMENT,	
					PRIMARY KEY (id)
				)  ENGINE=INNODB;`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"Table tasks is absent in schema1_\\d+.127.0.0.1:33060",
			},
		},
		"Schema Migrations": {
			schema1: []string{
				`CREATE TABLE IF NOT EXISTS schema_migrations (
					version VARCHAR(255) NOT NULL,
					UNIQUE KEY version_key(version)
				)  ENGINE=INNODB;`,
				`INSERT INTO schema_migrations values (20190815193300);`,
				`INSERT INTO schema_migrations values (20190817000000);`,
			},

			schema2: []string{
				`CREATE TABLE IF NOT EXISTS schema_migrations (
					version VARCHAR(255) NOT NULL,
					UNIQUE KEY version_key(version)
				)  ENGINE=INNODB;`,
				`INSERT INTO schema_migrations values (20190815193300);`,
				`INSERT INTO schema_migrations values (20190816000000);`,
			},
			expected: []string{
				"Differences found \\(1\\)",
				"\t- Some migrations are missing:",
				"\t\t- 127.0.0.1:33060",
				"\t\t\t- 20190816000000",
				"\t\t- 127.0.0.1:33062",
				"\t\t\t- 20190817000000",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cf, _ := NewFormatter("compact")
			result := RunDiff(t, test.schema1, test.schema2, cf)
			for _, expected := range test.expected {
				Regexp(t, expected, result)
			}
		})
	}
}
