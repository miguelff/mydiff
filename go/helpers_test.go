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
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
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
}

// MySQLIntegrationTest provides helper methods to perform integration tests
// aimed at checking schema differences
type MySQLIntegrationTest struct {
	s1Conn, s2Conn *sql.DB
}

// NewMySQLIntegrationTest returns the address of a new value of MySQLIntegrationTest
func NewMySQLIntegrationTest() *MySQLIntegrationTest {
	return &MySQLIntegrationTest{}
}

// Setup will spawn two mysql servers using docker-compose those servers will
// opinionatedly listen in ports 33060 and 33062.
func (m *MySQLIntegrationTest) Setup(t *testing.T) {
	t.Helper()
	execCmd("make", "db_down", "db_up")
	m.s1Conn = connect(S1DSN)
	m.s2Conn = connect(S2DSN)
}

// LoadSchemas will load sql1 and sql2 into the two servers under randomly generated
// schema names. The names will be returned to refer to them later. This allows for
// parallel test execution.
func (m *MySQLIntegrationTest) LoadSchemas(t *testing.T, sql1, sql2 []string) (schema1, schema2 string) {
	t.Helper()
	ts := time.Now().UnixNano()
	schema1 = fmt.Sprintf("schema1_%d", ts)
	schema2 = fmt.Sprintf("schema2_%d", ts)

	_, err := m.s1Conn.Exec("CREATE DATABASE " + schema1)
	if err != nil {
		log.Panic(err)
	}
	_, err = m.s1Conn.Exec("USE " + schema1)
	if err != nil {
		log.Panic(err)
	}

	for _, sql := range sql1 {
		_, err = m.s1Conn.Exec(sql)
		if err != nil {
			log.Panic(err)
		}
	}

	_, err = m.s2Conn.Exec("CREATE DATABASE " + schema2)
	if err != nil {
		log.Panic(err)
	}
	_, err = m.s2Conn.Exec("USE " + schema2)
	if err != nil {
		log.Panic(err)
	}

	for _, sql := range sql2 {
		_, err = m.s2Conn.Exec(sql)
		if err != nil {
			log.Panic(err)
		}
	}

	return
}

// Teardown will stop the two mysql servers, thus deleting all its information.
func (m *MySQLIntegrationTest) Teardown(t *testing.T) {
	t.Helper()
	execCmd("make", "db_down")
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

// execCmd executes the given shell command
func execCmd(cmd ...string) {
	log.Info("Executing command: ", strings.Join(cmd, " "))
	_, filename, _, _ := runtime.Caller(1)
	dir := path.Join(path.Dir(filename), "../")

	c := exec.Command(cmd[0], cmd[1:]...)
	c.Dir = dir
	out, err := c.Output()

	if err != nil {
		log.Error(err)
	}
	if out != nil {
		log.Info(string(out))
	}
}

// connect returns a *sql.DB connection to the given DSN
func connect(DSN string) *sql.DB {
	conn, _ := sql.Open("mysql", DSN)
	for {
		err := conn.Ping()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	return conn
}
