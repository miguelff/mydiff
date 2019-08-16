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
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	. "github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	sql1 := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id INT AUTO_INCREMENT,
			title CHAR(255) NOT NULL,
			PRIMARY KEY (id)
		)  ENGINE=INNODB;`,
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) NOT NULL,
			UNIQUE KEY version_key(version)
		)  ENGINE=INNODB;`,
		`INSERT INTO schema_migrations values (20190815193300);`,
		`INSERT INTO schema_migrations values (20190817000000);`,
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
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) NOT NULL,
			UNIQUE KEY version_key(version)
		)  ENGINE=INNODB;`,
		`INSERT INTO schema_migrations values (20190815193300);`,
		`INSERT INTO schema_migrations values (20190816000000);`,
	}

	s1Name, s2Name := TestCluster.LoadSchemas(t, sql1, sql2)
	from := NewServer1Schema(s1Name)
	to := NewServer2Schema(s2Name)

	diff := NewDiff(DSN1, DSN2, from, to, true, "schema_migrations.version")
	objectDiffs := diff.Compute()
	Equal(t, 3, len(objectDiffs))

	alters := objectDiffs[0].(*TableDiff).AlterClauses()
	Equal(t, 3, len(alters)) // modify tasks.id To make it bigint, title To make it a varchar, and add onwer_id\

	create, _ := objectDiffs[1].(*TableDiff)
	Equal(t, "owners", create.To.Name)

	migrations, _ := objectDiffs[2].(*MigrationsDiff)
	Equal(t, 1, len(migrations.Missing1))
	Equal(t, "20190816000000", migrations.Missing1[0])
	Equal(t, 1, len(migrations.Missing1))
	Equal(t, "20190817000000", migrations.Missing2[0])
}

func TestDiff_DontDiffMigrations(t *testing.T) {
	sql1 := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) NOT NULL,
			UNIQUE KEY version_key(version)
		)  ENGINE=INNODB;`,
		`INSERT INTO schema_migrations values (20190815193300);`,
		`INSERT INTO schema_migrations values (20190817000000);`,
	}

	sql2 := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) NOT NULL,
			UNIQUE KEY version_key(version)
		)  ENGINE=INNODB;`,
		`INSERT INTO schema_migrations values (20190815193300);`,
		`INSERT INTO schema_migrations values (20190816000000);`,
	}

	s1Name, s2Name := TestCluster.LoadSchemas(t, sql1, sql2)
	from := NewServer1Schema(s1Name)
	to := NewServer2Schema(s2Name)

	diff := NewDiff(DSN1, DSN2, from, to, false, "")
	objectDiffs := diff.Compute()
	Equal(t, 0, len(objectDiffs))
}

func TestDiff_WrongMigrationsColumnDefinition(t *testing.T) {
	oldLevel := logrus.GetLevel()
	defer func() {
		logrus.SetOutput(os.Stdout)
		logrus.SetLevel(oldLevel)
	}()
	mockWriter := NewMockWriter()
	logrus.SetOutput(mockWriter)
	logrus.SetLevel(logrus.WarnLevel)

	sql1 := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) NOT NULL,
			UNIQUE KEY version_key(version)
		)  ENGINE=INNODB;`,
		`INSERT INTO schema_migrations values (20190815193300);`,
		`INSERT INTO schema_migrations values (20190817000000);`,
	}

	sql2 := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) NOT NULL,
			UNIQUE KEY version_key(version)
		)  ENGINE=INNODB;`,
		`INSERT INTO schema_migrations values (20190815193300);`,
		`INSERT INTO schema_migrations values (20190816000000);`,
	}

	s1Name, s2Name := TestCluster.LoadSchemas(t, sql1, sql2)
	from := NewServer1Schema(s1Name)
	to := NewServer2Schema(s2Name)

	diff := NewDiff(DSN1, DSN2, from, to, true, "missing_table.version")
	objectDiffs := diff.Compute()
	Equal(t, 0, len(objectDiffs))

	diff = NewDiff(DSN1, DSN2, from, to, true, "schema_migrations.missing_column")
	objectDiffs = diff.Compute()
	Equal(t, 0, len(objectDiffs))

	Regexp(t, "Cannot retrieve migrations from missing_table.version", mockWriter.Entries[0])
	Regexp(t, "Cannot retrieve migrations from schema_migrations.missing_column", mockWriter.Entries[2])
}
