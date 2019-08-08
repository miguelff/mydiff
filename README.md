# `mydiff`

Compute the differences between two MySQL schemas.

mydiff is an alternative to [mysqldiff](https://docs.oracle.com/cd/E17952_01/mysql-utilities-1.5-en/mysqldiff.html#option_mysqldiff_difftype) that can output schema differences as migration scripts for [gh-ost](https://github.com/github/gh-ost), [pt-ost](https://www.percona.com/doc/percona-toolkit/LATEST/pt-online-schema-change.html), or as [active record migrations](https://edgeguides.rubyonrails.org/active_record_migrations.html)

It's written in golang and can compute the differences in parallel, to speed up the comparision of large schemas.

## Usage

`mydiff --server1=user:pass@host:port:socket --server2=user:pass@host:port:socket db1.object1:db2.object1 db3:db4`

### Options

*`--version`* 

show program's version number and exit

*`--help`*

display a help message and exit

*`--server1=SERVER1`*

connection information for first server in the form: `<user>[:<password>]@<host>[:<port>][:<socket>]` or `<login-path>[:<port>][:<socket>]` or path yo `my.cnf` file.

*`--server2=SERVER2`*

connection information for first server in the form: `<user>[:<password>]@<host>[:<port>][:<socket>]` or `<login-path>[:<port>][:<socket>]` or path yo `my.cnf` file.

*`-d DIFFTYPE`, `--difftype=DIFFTYPE`*

display differences in context format in one of the following formats: [`sql|gh-ost|pt-osc|activerecord`] (default: `sql`).

## Testing `mysqldiff`

* `make test` will run [golangci-lint](https://github.com/golangci/golangci-lint) and unit tests.

* `make demo` will use docker compose to:
    - Spawn two MySQL servers, and load slightly different schemas dumps.
    - Run mysqldiff in another container, computing the differences between the two server schemas, using
    different command options.
    
## Design decisions and trade-offs

TBD

## Missing features

TBD

## License

`mydiff` is licensed under the [MIT license](https://github.com/miguelff/mydiff/blob/master/LICENSE) 

## Authors

`mydiff` is designed, authored, reviewed and tested by [@miguelff](https://github.com/miguelff)
