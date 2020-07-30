[codecov]: https://codecov.io/gh/magicalbanana/dbal
[codecov badge master]: https://codecov.io/gh/magicalbanana/dbal/branch/master/graph/badge.svg
[drone]: https://cloud.drone.io
[drone repo]: https://cloud.drone.io/magicalbanana/dbal
[drone badge master]: https://cloud.drone.io/api/badges/magicalbanana/dbal/status.svg
[doc]: https://godoc.org/github.com/magicalbanana/dbal
[doc badge]: https://godoc.org/github.com/magicalbanana/dbal?status.svg
[go report card]: https://goreportcard.com/report/github.com/magicalbanana/dbal
[go report card badge]: https://goreportcard.com/badge/github.com/magicalbanana/dbal

| master                                                    |
| -                                                         |
| [![drone repo][drone badge master]][drone repo]           |
| [![codecov][codecov badge master]][codecov]               |
| [![doc][doc badge]][doc]                                  |
| [![go report card][go report card badge]][go report card] |

# dbal - DB Access Layer

## Description

This package wraps the `database/sql` package by allowing you to specify
parameters in a map and converting the SQL parameters to positional parameters
to leverage the SQL santizaton provided by the `database/sql` package.

## Usage

TBD ...

## TODO

- [ ] add support for transactions
- [ ] benchmark SQL parsing
