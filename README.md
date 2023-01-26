# becrypt
**Generate and test bcrypt hashes from a CLI**

[![GitHub](https://img.shields.io/github/license/pepa65/becrypt.svg)](LICENSE)
[![run-ci](https://github.com/pepa65/becrypt/actions/workflows/ci.yml/badge.svg)](https://github.com/pepa65/becrypt/actions/workflows/ci.yml) 
* Version: 1.2.3
* License: [MIT](LICENSE)
* Repo: `github.com/pepa65/becrypt`
* Modified interface from `github.com/shoenig/bcrypt-tool`:
  - Shorter & simpler, and only a command for the least used option.
  - No password on the commandline (either piped-in or asked for interactively).
  - Functionally compatible (both use `golang.org/x/crypto/bcrypt` under the hood).

## Usage
```
becrypt v1.2.3 - Generate and test bcrypt hashes from a CLI
Repo:   github.com/pepa65/becrypt
Usage:  becrypt HASH | TEST | COST | HELP
 HASH:  [<cost>]:               Generate a hash from the password
                                (optional <cost>: 4..31, default: 10)
 TEST:  <hash>:                 Test the password against <hash>
 COST:  cost|-c|--cost <hash>:  Display the cost of <hash>
 HELP:  help|-h|--help:         Display this help text
The password can be piped-in or prompted for, is cut off after 72 characters.
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

### Generate Hash from a Password
```bash
becrypt  # A password will be asked for interactively

printf 'p4ssw0rd' |becrypt
```

### Generate Hash from a Password with given Cost
```bash
becrypt 31  # A password will be asked for interactively

printf 'p4ssw0rd' |becrypt 4
```

### Determine if Password matches Hash
```bash
# A password will be asked for interactively
becrypt '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'

printf 'p4ssw0rd' |becrypt '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'
```

### Determine processing Cost of Hash
```bash
becrypt cost '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'
```

