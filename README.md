# becrypt
**CLI tool for generating and checking bcrypt hashes**

* License: [MIT](LICENSE) [![GitHub](https://img.shields.io/github/license/pepa65/becrypt.svg)](LICENSE)
* Repo: `github.com/pepa65/becrypt`
* Modified interface from `github.com/shoenig/bcrypt-tool`
  - Shorter & simpler, and only a command for the least used option
  - No password on the commandline (either piped-in or asked for interactively)

## Usage
```
becrypt [<cost>] | <hash> | cost <hash>
    becrypt [<cost>]:     Generate a hash from the piped-in or entered password
                          (optional <cost>: 4..31, default: 10)
    becrypt <hash>:       Check the piped-in or entered password against <hash>
    becrypt cost <hash>:  Display the cost of <hash>
```

## Install from Releases

* The `becrypt` tool is available from the [Releases](https://github.com/pepa65/becrypt/releases) page.
* Pre-compiled for:
  - Linux
  - Windows
  - MacOS
  - FreeBSD
  - OpenBSD
  - Plan9

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

### Determine Cost of Hash
```bash
becrypt cost '$2a$10$nWFwjoFo4zhyVosdYMb6XOxZqlVB9Bk0TzOvmuo16oIwMZJXkpanW'
```

