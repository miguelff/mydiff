# `mydiff`

Compute the differences between two MySQL schemas.

mydiff is an alternative to [mysqldiff](https://docs.oracle.com/cd/E17952_01/mysql-utilities-1.5-en/mysqldiff.html#option_mysqldiff_difftype) that can output schema differences as migration scripts for [gh-ost](https://github.com/github/gh-ost), or as [active record migrations](https://edgeguides.rubyonrails.org/active_record_migrations.html)

It's written in golang and can compute the differences in parallel, to speed up the comparision of large schemas.

## Usage (`mydiff --help`)

```
NAME:
   mydiff - Compute the differences between two MySQL schemas

USAGE:
   mydiff --server1=user:pass@host:port:socket --server2=user:pass@host:port:socket schema_name [schema_name in server2]

GLOBAL OPTIONS:
   --server1 value             connection information for second server in the form of a DSN (<user>[:<password>]@<host>[:<port>][:<socket>]) or path to socket file.
   --server2 value             connection information for second server in the form of a DSN (<user>[:<password>]@<host>[:<port>][:<socket>]) or path to socket file.
   -d value, --difftype value  display differences in one of the following formats: [sql|gh-ost|ar] (default: "sql")
   -r, --reverse               show diff in reverse direction, from server2 to server1
   -v, --version               display version
   -h, --help                  display this help

COPYRIGHT:
   Copyright 2019 Miguel Fernández. Licensed under MIT license
```

## Testing `mysqldiff`

* `make test` will run [golangci-lint](https://github.com/golangci/golangci-lint) and unit tests.

* `make demo` will use docker compose to:
    - Spawn two MySQL servers, and load slightly different schemas dumps.
    - Run mysqldiff in another container, computing the differences between the two server schemas, using
    different command options.
    
## Design decisions and trade-offs

TBD

## Missing features

- [ ] Output diff in the different formats (ar migrations, gh-ost, etc)
- [ ] parse server connection descriptors in different formats that are more flexible than golang sqlx parser
- [ ] demo
- [ ] unit tests

## License

`mydiff` is licensed under the [MIT license](https://github.com/miguelff/mydiff/blob/master/LICENSE) 

## Authors

`mydiff` is designed, authored, reviewed and tested by [@miguelff](https://github.com/miguelff)
