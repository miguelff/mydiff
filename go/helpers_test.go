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
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/skeema/tengo"

	log "github.com/sirupsen/logrus"
)

const (
	S1DSN = "root@tcp(127.0.0.1:33060)/"
	S2DSN = "root@tcp(127.0.0.1:33062)/"
)

var TestCluster *MySQLCluster

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)

	// set log level from env (1: Fatal, 5: debug)
	level := os.Getenv("LOG_LEVEL")
	i, err := strconv.ParseUint(level, 10, 32)
	if err != nil {
		i = uint64(log.WarnLevel)
	}
	log.SetLevel(log.Level(i))
	mysql.SetLogger(log.StandardLogger())

	TestCluster = &MySQLCluster{
		s1: connect(S1DSN, 5*time.Second),
		s2: connect(S2DSN, 5*time.Second),
	}
}

// MySQLCluster holds the connections to two servers
type MySQLCluster struct {
	s1, s2 *sql.DB
}

// LoadSchemas will load sql1 and sql2 into the two servers under randomly generated
// schema names. The names will be returned to refer to them later. This allows for
// parallel test execution.
func (m *MySQLCluster) LoadSchemas(t *testing.T, sql1, sql2 []string) (schema1, schema2 string) {
	t.Helper()
	ts := time.Now().UnixNano()
	schema1 = fmt.Sprintf("schema1_%d", ts)
	schema2 = fmt.Sprintf("schema2_%d", ts)

	_, err := m.s1.Exec("CREATE DATABASE " + schema1)
	if err != nil {
		log.Panic(err)
	}
	_, err = m.s1.Exec("USE " + schema1)
	if err != nil {
		log.Panic(err)
	}

	for _, sql := range sql1 {
		_, err = m.s1.Exec(sql)
		if err != nil {
			log.Panic(err)
		}
	}

	_, err = m.s2.Exec("CREATE DATABASE " + schema2)
	if err != nil {
		log.Panic(err)
	}
	_, err = m.s2.Exec("USE " + schema2)
	if err != nil {
		log.Panic(err)
	}

	for _, sql := range sql2 {
		_, err = m.s2.Exec(sql)
		if err != nil {
			log.Panic(err)
		}
	}

	return
}

// NewServer1Schema returns a tengo.Schema value denoted by the given name in the server1
func NewServer1Schema(name string) *tengo.Schema {
	return newSchema(S1DSN, name)
}

// NewServer2Schema returns a tengo.Schema value denoted by the given name in the server2
func NewServer2Schema(name string) *tengo.Schema {
	return newSchema(S2DSN, name)
}

// newSchema returns the address of a new tengo.Schema described by the given
// DSN and schema names
func newSchema(DSN, schema string) *tengo.Schema {
	i, err := tengo.NewInstance("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	s, err := i.Schema(schema)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

// connect returns a *sql.DB connection to the given DSN
func connect(DSN string, timeout time.Duration) *sql.DB {
	ticker := time.NewTicker(1 * time.Second)
	timer := time.After(timeout)
	conn, _ := sql.Open("mysql", DSN)

	for {
		select {
		case <-ticker.C:
			err := conn.Ping()
			if err == nil {
				return conn
			}
		case <-timer:
			log.Fatalf("Timeout trying to connect to %s after %d seconds. Forgot to run `make db_up`?", DSN, timeout/1e9)
			os.Exit(1)
		}
	}

}
