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
	app.UsageText = "mydiff --server1=user:pass@tcp(host:port)/ --server2=user:pass@tcp(host:port)/ GLOBAL OPTIONS schema_name"

	app.HideHelp = true
	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server1",
			Usage: "connection information for second server in the form of a DSN (<user>[:<password>]@tcp(<host>[:<port>])/)",
		},
		cli.StringFlag{
			Name:  "server2",
			Usage: "connection information for second server in the form of a DSN (<user>[:<password>]@tcp(<host>[:<port>])/)",
		},
		cli.StringFlag{
			Name:  "d, diff-type",
			Value: "compact",
			Usage: "display differences in one of the following formats: [sql|compact]",
		},
		cli.BoolFlag{
			Name:  "diff-migrations",
			Usage: "if the schema has a migrations table, compute its difference. Works only with compact formatting",
		},
		cli.StringFlag{
			Name:  "diff-migrations-column",
			Value: "schema_migrations.version",
			Usage: "if --diff-migrations is enabled, this flag will determine which column values to compare in both schemas",
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
		server1, err := tengo.NewInstance(driver, mydiff.ParseDSN(c.GlobalString("server1")).FormatDSN())
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("server1 has to be a server DSN. Error: %s", err.Error()), EServInvalid)
		}
		server2, err := tengo.NewInstance(driver, mydiff.ParseDSN(c.GlobalString("server2")).FormatDSN())
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

		formatter, err := mydiff.NewFormatter(c.GlobalString("diff-type"))
		if err != nil {
			return cli.NewExitError(err, EUnkownFormatter)
		}
		includeMigrations := c.GlobalBool("diff-migrations")
		migrationsCol := c.GlobalString("diff-migrations-column")

		if c.GlobalBool("reverse") {
			tmp := to
			to = from
			from = tmp

			sTmp := server1
			server1 = server2
			server2 = sTmp
		}

		diff := mydiff.NewDiff(server1.BaseDSN, server2.BaseDSN, from, to, includeMigrations, migrationsCol)
		result := formatter.Format(diff)
		fmt.Print(result)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
