> command `prunedb` will drop mysql databases that being with the `test_`
> prefix.

[![Go
Reference](https://pkg.go.dev/badge/github.com/wayneashleyberry/prunedb.svg)](https://pkg.go.dev/github.com/wayneashleyberry/prunedb)
![build](https://github.com/wayneashleyberry/prunedb/workflows/build/badge.svg)
[![Go Report
Card](https://goreportcard.com/badge/github.com/wayneashleyberry/prunedb)](https://goreportcard.com/report/github.com/wayneashleyberry/prunedb)

### Installation

```sh
go get github.com/wayneashleyberry/prunedb
```

### Configuration

This application is configured via the environment. The following environment variables can be used:

```sh
KEY                TYPE      DEFAULT      REQUIRED    DESCRIPTION
TESTDB_HOST        String    127.0.0.1
TESTDB_PORT        String    3306
TESTDB_USER        String                 true
TESTDB_PASSWORD    String                 true
```

### Usage

```sh
Drop lingering test databases

Usage:
  prunedb [flags]

Flags:
      --commit   commit query
  -h, --help     help for prunedb
```
