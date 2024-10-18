[![Go Report Card](https://goreportcard.com/badge/github.com/pepa65/becrypt)](https://goreportcard.com/report/github.com/pepa65/becrypt)
[![GoDoc](https://godoc.org/github.com/pepa65/becrypt?status.svg)](https://godoc.org/github.com/pepa65/becrypt)
[![GitHub](https://img.shields.io/github/license/pepa65/becrypt.svg)](LICENSE)
[![run-ci](https://github.com/pepa65/becrypt/actions/workflows/ci.yml/badge.svg)](https://github.com/pepa65/becrypt/actions/workflows/ci.yml) 

# becrypt
**Generate and check bcrypt hashes from a CLI**

* Version: 1.4.1
* License: [MIT](LICENSE)
* Repo: `github.com/pepa65/becrypt`
* Modified interface from `github.com/shoenig/bcrypt-tool`:
  - Shorter & simpler, and only a command for the least used option.
  - No password on the commandline (either piped-in or asked for interactively).
  - Functionally compatible (both use `golang.org/x/crypto/bcrypt` under the hood).

## Usage
```
becrypt v1.4.1 - Generate and check bcrypt hashes from a CLI
Repo:   github.com/pepa65/becrypt
Usage:  becrypt OPTION
    Options:
        help|-h|--help           Display this HELP text.
        cost|-c|--cost <hash>    Display the COST of bcrypt <hash>.
        <hash> [-q|--quiet]      CHECK the password(^) against bcrypt <hash>.
        [<cost>]                 Generate a HASH from the given password(^).
                                 (Optional <cost>: 4..31, default: 10.)
(^) Password: can be piped-in or prompted for, it gets cut off after 72 bytes.
    Longer ones are accepted without warning, using only the first 72 bytes!
```

## Install from Releases

* The `becrypt` tool is available from the [Releases](https://github.com/pepa65/becrypt/releases) page.
* Pre-compiled for:
  - OSX amd64 arm64
  - Linux amd64 386 arm64
  - FreeBSD amd64 386 arm64
  - Openbsd amd64 386 arm64
  - Windows amd64 386 arm64
  - Plan9 amd64 386

## Build from source with Go
```bash
go get github.com/pepa65/becrypt
```

## Examples
**Quote the password/hash! (Depending on your shell.)**

### COST: Determine processing Cost of Hash
```bash
becrypt cost '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'
```

The result of a COST command is a plaintext 10-based number on stdout with returncode 0,
unless the hash is malformed, then an error results for a returncode bigger than 0).

### CHECK: Determine if Password matches Hash
```bash
# A password will be asked for interactively
becrypt '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'

printf 'p4ssw0rd' |becrypt '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'
```

The result of a CHECK command is a plaintext 'Y' or 'N' on stdout,
with corresponding returncodes 0 and 1.
If the `-q` or `--quiet` flag is given, no stdout is produced, only the returncode.

### HASH: Generate Hash from a Password
```bash
becrypt  # A password will be asked for interactively

printf 'p4ssw0rd' |becrypt
```

The result of a HASH command is the hash on stdout, with a returncode of 0.

### HASH: Generate Hash from a Password with given Cost
```bash
becrypt 31  # A password will be asked for interactively

printf 'p4ssw0rd' |becrypt 4
```

The processing cost scales exponentially with 2^cost,
so a cost increase of 1 doubles the processing time needed.
So higher cost numbers will take a while!

## Release management
* Change version in `README.md` (2 places) and `main.go`.
* Commit changes
* Add tag: `git tag -a <tag>`
* Push commit(s)
* Do: `goreleaser --clean` (goreleaser binary at: https://github.com/goreleaser/goreleaser/releases)
