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
	DSN1        = "root@tcp(127.0.0.1:33060)/"
	DSN2        = "root@tcp(127.0.0.1:33062)/"
	connTimeout = 60 * time.Second
)

var TestCluster *MySQLCluster

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)

	// set log level From env (1: Fatal, 5: debug)
	level := os.Getenv("LOG_LEVEL")
	i, err := strconv.ParseUint(level, 10, 32)
	if err != nil {
		i = uint64(log.ErrorLevel)
	}
	log.SetLevel(log.Level(i))
	_ = mysql.SetLogger(log.StandardLogger())

	TestCluster = &MySQLCluster{
		s1: connect(DSN1, connTimeout),
		s2: connect(DSN2, connTimeout),
	}
}

// MySQLCluster holds the connections To two servers
type MySQLCluster struct {
	s1, s2 *sql.DB
}

// LoadSchemas will load sql1 and sql2 into the two servers under randomly generated
// schema names. The names will be returned To refer To them later. This allows for
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
	return newSchema(DSN1, name)
}

// NewServer2Schema returns a tengo.Schema value denoted by the given name in the server2
func NewServer2Schema(name string) *tengo.Schema {
	return newSchema(DSN2, name)
}

// RunDiff runs a diff between the two given schemas, applying the formatter and format options also given
func RunDiff(t *testing.T, schema1 []string, schema2 []string, formatter Formatter) interface{} {
	t.Helper()
	s1Name, s2Name := TestCluster.LoadSchemas(t, schema1, schema2)
	from := NewServer1Schema(s1Name)
	to := NewServer2Schema(s2Name)
	diff := NewDiff(DSN1, DSN2, from, to, true, "schema_migrations.version")
	return formatter.Format(diff)
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

// connect returns a *sql.DB connection To the given DSN
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
			log.Fatalf("Timeout trying To connect To %s after %d seconds. Forgot To run `make db_up`?", DSN, timeout/1e9)
			os.Exit(1)
		}
	}
}

// Mock writer is a writer that stores the
// written bytes as strings in its Entries fields
type MockWriter struct {
	Entries []string
}

// NewMockWriter returns the address of a new MockWriter value
func NewMockWriter() *MockWriter {
	return &MockWriter{Entries: []string{}}
}

// Write stores p as a string in Entries
func (w *MockWriter) Write(p []byte) (int, error) {
	w.Entries = append(w.Entries, string(p))
	return len(p), nil
}
