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

package main

import (
	"fmt"
	"log"
	"os"

	mydiff "github.com/miguelff/mydiff/go"

	"github.com/skeema/tengo"
	"github.com/urfave/cli"
)

const driver = "mysql"

const (
	ESchemaNameNotProvided = iota + 1
	EServInvalid
	EMissingSchema
	EUnkownFormatter
)

func main() {
	app := cli.NewApp()
	app.Name = "mydiff"
	app.Version = mydiff.Version
	app.Usage = "Compute the differences between two MySQL schemas"
	app.Copyright = "Copyright 2019 Miguel Fernández. Licensed under MIT license"
	app.UsageText = "mydiff --server1=user:pass@host:port:socket --server2=user:pass@host:port:socket schema_name [schema_name in server2]"

	app.HideHelp = true
	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server1",
			Usage: "connection information for second server in the form of a DSN (<user>[:<password>]@<host>[:<port>][:<socket>]) or path to socket file.",
		},
		cli.StringFlag{
			Name:  "server2",
			Usage: "connection information for second server in the form of a DSN (<user>[:<password>]@<host>[:<port>][:<socket>]) or path to socket file.",
		},
		cli.StringFlag{
			Name:  "d, difftype",
			Value: "sql",
			Usage: "display differences in one of the following formats: [sql|ghost|ar]",
		},
		cli.BoolFlag{
			Name:  "r, reverse",
			Usage: "show diff in reverse direction, from server2 to server1",
		},
		cli.BoolFlag{
			Name:  "v, version",
			Usage: "display version",
		},
		cli.BoolFlag{
			Name:  "h, help",
			Usage: "display this help",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.GlobalBool("help") {
			return cli.ShowAppHelp(c)
		}

		if c.GlobalBool("version") {
			cli.ShowVersion(c)
			return nil
		}

		schema1 := c.Args().Get(0)
		schema2 := c.Args().Get(1)

		if schema2 == "" {
			schema2 = schema1
		}
		if schema1 == "" {
			return cli.NewExitError("schema_name has to be provided", ESchemaNameNotProvided)
		}
		server1, err := tengo.NewInstance(driver, c.GlobalString("server1"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("server1 has to be a server DSN. Error: %s", err.Error()), EServInvalid)
		}
		server2, err := tengo.NewInstance(driver, c.GlobalString("server2"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("server2 has to be a server DSN. Error: %s", err.Error()), EServInvalid)
		}
		from, err := server1.Schema(schema1)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("server1 doesn't contain schema %s. Error: %s", schema1, err.Error()), EMissingSchema)
		}
		to, err := server2.Schema(schema2)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("server2 doesn't contain schema %s. Error: %s", schema2, err.Error()), EMissingSchema)
		}

		formatter, err := mydiff.NewFormatter(c.GlobalString("difftype"))
		if err != nil {
			return cli.NewExitError(err, EUnkownFormatter)
		}

		if c.GlobalBool("reverse") {
			tmp := to
			to = from
			from = tmp
		}

		diff := mydiff.NewDiff(from, to)
		result := formatter.Format(diff)
		fmt.Print(result)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
