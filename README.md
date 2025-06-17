[![Go Report Card](https://goreportcard.com/badge/github.com/pepa65/becrypt)](https://goreportcard.com/report/github.com/pepa65/becrypt)
[![GoDoc](https://godoc.org/github.com/pepa65/becrypt?status.svg)](https://godoc.org/github.com/pepa65/becrypt)
[![GitHub](https://img.shields.io/github/license/pepa65/becrypt.svg)](LICENSE)
[![run-ci](https://github.com/pepa65/becrypt/actions/workflows/ci.yml/badge.svg)](https://github.com/pepa65/becrypt/actions/workflows/ci.yml) 

# becrypt
**Generate and check bcrypt hashes from a CLI**

* Version: 1.5.2
* License: [MIT](LICENSE)
* Repo: `github.com/pepa65/becrypt`
* Modified interface from `github.com/shoenig/bcrypt-tool`:
  - Shorter & simpler, and only a command for the least used option.
  - No password on the commandline (it either gets piped-in or is asked for interactively).
    Any final newline gets cut off, so `echo 'pw' |becrypt` and `becrypt <<<'pw' can be used.
  - Functionally compatible (both use `golang.org/x/crypto/bcrypt` under the hood).

## Usage
```
becrypt v1.5.2 - Generate and check bcrypt hashes from a CLI
Repo:   github.com/pepa65/becrypt
Usage:  becrypt OPTION
    Options:
        help|-h|--help           Display this HELP text.
        cost|-c|--cost <hash>    Display the COST of bcrypt <hash>.
        <hash> [-q|--quiet]      CHECK the password(^) against bcrypt <hash>.
        [<cost>]                 Generate a HASH from the given password(^).
                                 (Optional <cost>: 4..31, default: 10.)
(^) Password: can piped-in or prompted for, a final newline will get cut off.
    Passwords longer than 72 bytes are accepted & get cut off without warning.
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
```
becrypt cost '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'
```

Output:
```
10
```

The result of a COST command is a plaintext 10-based number on stdout with returncode 0,
unless the hash is malformed, then an error results for a returncode bigger than 0).

### CHECK: Determine if Password matches Hash
```
# A password will be asked for interactively
becrypt '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'

echo 'p4ssw0rd' |becrypt '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'
```

The result of a CHECK command is a plaintext 'true' or 'false' on stdout,
with corresponding returncodes 0 and 1.
If the `-q` or `--quiet` flag is given, no stdout is produced, only the returncode.

### HASH: Generate Hash from a Password
```
# A password will be asked for interactively
becrypt

echo 'p4ssw0rd' |becrypt

becrypt <<<'p4ssw0rd'
```

The result of a HASH command is the hash on stdout, with a returncode of 0.

### HASH: Generate Hash from a Password with given Cost
```bash
becrypt 31  # A password will be asked for interactively

printf 'p4ssw0rd' |becrypt 4

becrypt 12 <<<'p4ssw0rd'
```

The hashing processing cost scales exponentially with 2^cost,
so a cost increase of 1 doubles the processing time needed.
So higher cost numbers will take significantly longer,
a cost increase of 10 takes more than a 1000 times longer!

## Release management
* Change version in `README.md` (2 places) and `main.go`.
* Commit changes
* Add tag: `git tag -a <tag>`
* Push commit(s)
* Do: `goreleaser --clean` (goreleaser binary at: https://github.com/goreleaser/goreleaser/releases)
