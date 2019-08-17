# `mydiff`

<!-- Uncomment when opensourced
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/miguelff/mydiff)
-->

[![Travis](https://travis-ci.com/miguelff/mydiff.svg?token=1bFaTyv8B89uBs2sxt7M&branch=master)](https://travis-ci.org/miguelff/mydiff)
[![Coverage Status](https://coveralls.io/repos/github/miguelff/mydiff/badge.svg?branch=master&t=41u1ce)](https://coveralls.io/github/miguelff/mydiff?branch=master)

<img width="1024" alt="mydiff logo" src="https://user-images.githubusercontent.com/210307/62741731-4f22de00-ba3c-11e9-89ee-da12f92e0b4f.png">

Compute the differences between two MySQL schemas.

`mydiff` is an alternative to [mysqldiff](https://docs.oracle.com/cd/E17952_01/mysql-utilities-1.5-en/mysqldiff.html#option_mysqldiff_difftype) written in golang as a thin wrapper of [skeema/tengo](github.com/skeema/tengo/)

## Usage (`mydiff --help`)

```
NAME:
   mydiff - Compute the differences between two MySQL schemas

USAGE:
   mydiff --server1=user:pass@tcp(host:port)/ --server2=user:pass@tcp(host:port)/ GLOBAL OPTIONS schema_name

GLOBAL OPTIONS:
   --server1 value                 connection information for second server in the form of a DSN (<user>[:<password>]@tcp(<host>[:<port>]))
   --server2 value                 connection information for second server in the form of a DSN (<user>[:<password>]@tcp(<host>[:<port>]))
   -d value, --diff-type value     display differences in one of the following formats: [sql|compact] (default: "compact")
   --diff-migrations               if the schema has a migrations table, compute its difference. Works only with compact formatting
   --diff-migrations-column value  if --diff-migrations is enabled, this flag will determine which column values to compare in both schemas (default: "schema_migrations.version")
   -r, --reverse                   show diff in reverse direction, from server2 to server1
   -v, --version                   display version
   -h, --help                      display this help

COPYRIGHT:
   Copyright 2019 Miguel Fernández. Licensed under MIT license
```

## Installation

`make build` build will generate in `.build/mydiff` a linux binary with all the dependencies statically linked. The binary will be ready to be used inside any docker image or native linux distribution.

## Testing

* `make test` will run [golangci-lint](https://github.com/golangci/golangci-lint) and integration tests.

* `make demo` runs an interactive demo using `mydiff` to compute schema differences using a variety of command options.

Both the tests and the demo will use docker-compose to spawn two mysql servers and thus mimic a real usage scenario.
In addition, you expect to have a ruby interpreter in your system to run the demo.
    
## Design decisions and trade-offs

[skeema/tengo](github.com/skeema/tengo/) does a great job at computing differences in schemas, thus this tool does not 
reinvent the wheel and instead provides a thin wrapper atop tengo, which focuses on providing a replacement for [mysqldiff](https://docs.oracle.com/cd/E17952_01/mysql-utilities-1.5-en/mysqldiff.html#option_mysqldiff_difftype)
that is straightly usable from the shell, and specialized on computing schema changes, but simpler and easier to use than the more complete [skeema/skeema](https://github.com/skeema/skeema). 

Also, because [skeema/tengo](github.com/skeema/tengo/) provides [tests to ensure diffs are computed correctly](https://github.com/skeema/tengo/blob/master/diff_test.go), `mysqldiff` tests focuses instead on:
 - Ensuring `tengo.Diff` objects are parsed correctly to the different output formats. 
 - Ensuring the different formatters generate correct output.  

## Limitations and Missing features

- [ ] Detecting changes in auto-increment initial values is not supported. This can be implemented by querying the auto increment values on the server directly, however this was not an initial requirement for the project and thus is left out of the scope of this first version of the tool.
- [ ] Changes in encoding are detected, however the formatter only displays the encoding in the second schema being compared as tengo loses information about how it was before. This can be fixed by querying the DB on server1 and inspecting the table collation and encoding, but this is left out of the scope as the compact output informs about a mismatch in encoding pretty clearly. 

## License

`mydiff` is licensed under the [MIT license](https://github.com/miguelff/mydiff/blob/master/LICENSE), and it uses [third-party libraries](https://github.com/miguelff/mydiff/blob/master/go.mod) that have their own licenses.

## Authors

`mydiff` is authored by [@miguelff](https://github.com/miguelff)
