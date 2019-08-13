# `mydiff`

<img width="1024" alt="mydiff logo" src="https://user-images.githubusercontent.com/210307/62741731-4f22de00-ba3c-11e9-89ee-da12f92e0b4f.png">

Compute the differences between two MySQL schemas.

`mydiff` is an alternative to [mysqldiff](https://docs.oracle.com/cd/E17952_01/mysql-utilities-1.5-en/mysqldiff.html#option_mysqldiff_difftype) written in golang as a thin wrapper of [skeema/tengo](github.com/skeema/tengo/)

## Usage (`mydiff --help`)

```
NAME:
   mydiff - Compute the differences between two MySQL schemas

USAGE:
   mydiff --server1=user:pass@host:port:socket --server2=user:pass@host:port:socket schema_name [schema_name in server2]

GLOBAL OPTIONS:
   --server1 value              connection information for second server in the form of a DSN (<user>[:<password>]@<host>[:<port>][:<socket>]) or path to socket file.
   --server2 value              connection information for second server in the form of a DSN (<user>[:<password>]@<host>[:<port>][:<socket>]) or path to socket file.
   -d value, --diff-type value  display differences in one of the following formats: [sql|compact] (default: "sql")
   --diff-opts value            options to pass through to the different diff-type formatters
   -r, --reverse                show diff in reverse direction, from server2 to server1
   -v, --version                display version
   -h, --help                   display this help

COPYRIGHT:
   Copyright 2019 Miguel Fernández. Licensed under MIT license
```

## Installation

`make build` build will generate in `.build/mydiff` a linux binary with all the dependencies statically linked. The binary will be ready to be used inside any docker image or native linux distribution.

## Testing

* `make test` will run [golangci-lint](https://github.com/golangci/golangci-lint) and integration tests.

* `make demo` load slightly different schemas dumps in two different servers and use `mydiff` the differences using the different command options.

Both the tests and the demo will use docker-compose two spawn two mysql servers and thus mimic a real usage scenario.
    
## Design decisions and trade-offs

[skeema/tengo](github.com/skeema/tengo/) does a great job at computing differences in schemas, thus this tool does not 
reinvent the wheel and instead provides a thin wrapper atop tengo, which focuses on providing a replacement for [mysqldiff](https://docs.oracle.com/cd/E17952_01/mysql-utilities-1.5-en/mysqldiff.html#option_mysqldiff_difftype)
that is straightly usable from the shell, and specialized on computing schema changes, but simpler and easier to use than the more complete [skeema/skeema](https://github.com/skeema/skeema). 

Also, because [skeema/tengo](github.com/skeema/tengo/) provides [tests to ensure diffs are computed correctly](https://github.com/skeema/tengo/blob/master/diff_test.go), `mysqldiff` tests focuses instead on:
 - Ensuring `tengo.Diff` objects are parsed correctly to the different output formats. 
 - Ensuring the different formatters generate correct output.  

## Missing features

- [ ] parse server connection descriptors in different formats that are more flexible than golang sqlx parser
- [ ] demo
- [ ] published releases on GitHub

## License

`mydiff` is licensed under the [MIT license](https://github.com/miguelff/mydiff/blob/master/LICENSE), and it uses [third-party libraries](https://github.com/miguelff/mydiff/blob/master/go.mod) that have their own licenses.

## Authors

`mydiff` is authored by [@miguelff](https://github.com/miguelff)
