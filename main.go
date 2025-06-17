// becrypt - Generate and check bcrypt hashes from a CLI

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

const (
	version  = "1.5.2"
	pwMaxLen = 72
	nl = 10
)

var (
	self  = ""
	quiet = false
)

func showHelp(msg string) { // I:self,version,pwMaxLen
	helptxt := self + " v" + version + ` - Generate and check bcrypt hashes from the CLI
Repo:   github.com/pepa65/becrypt
Usage:  ` + self + ` OPTION
    Options:
        help|-h|--help           Display this HELP text.
        cost|-c|--cost <hash>    Display the COST of bcrypt <hash>.
        <hash> [-q|--quiet]      CHECK the password(^) against bcrypt <hash>.
        [<cost>]                 Generate a HASH from the given password(^).
                                 (Optional <cost>: ` +
		strconv.Itoa(bcrypt.MinCost) + ".." + strconv.Itoa(bcrypt.MaxCost) +
		", default: " + strconv.Itoa(bcrypt.DefaultCost) + `.)
(^) Password: can piped-in or prompted for, a final newline will get cut off.
    Passwords longer than ` +
		strconv.Itoa(pwMaxLen) + " bytes are accepted & get cut off without warning."
	fmt.Fprintln(os.Stderr, helptxt)
	if msg != "" {
		fmt.Fprintln(os.Stderr, "Abort: "+msg)
		os.Exit(2)
	}
	os.Exit(0)
}

func main() { // IO:self
	hash, cmd, cost := "", "", bcrypt.DefaultCost
	for _, arg := range os.Args {
		if self == "" {
			selves := strings.Split(arg, "/")
			self = selves[len(selves)-1]
			continue
		}
		if cmd == "COST" {
			if hash != "" {
				showHelp("Too many arguments for cost: " + arg)
			}
			hash = arg
			continue
		}
		if cmd == "CHECK" {
			if arg == "-q" || arg == "--quiet" {
				quiet = true
				continue
			}
			if hash != "" {
				showHelp("Hash already given, too many arguments: " + arg)
			} else {
				hash = arg
				continue
			}
		}
		if cmd == "" && (arg == "cost" || arg == "-c" || arg == "--cost") {
			cmd = "COST"
			continue
		}
		if arg == "help" || arg == "-h" || arg == "--help" {
			showHelp("")
		}
		if cmd == "" {
			c, e := strconv.Atoi(arg)
			if e == nil { // Integer: a hash command
				cmd, cost = "HASH", c
				continue
			}
			if arg == "-q" || arg == "--quiet" {
				cmd, quiet = "CHECK", true
				continue
			}
			cmd, hash = "CHECK", arg
		} else {
			showHelp("Too many arguments: " + arg)
		}
	}
	if cmd == "" {
		cmd = "HASH"
	}
	if hash != "" { // Check hash
		hashpart := strings.Split(hash, "$")
		switch true {
		case len(hashpart) != 4:
			fmt.Fprintln(os.Stderr, "Abort: Exactly 3 '$' in a bcrypt hash, invalid hash: "+hash)
			os.Exit(3)
		case len(hashpart[0]) > 0:
			fmt.Fprintln(os.Stderr, "Abort: A proper bcrypt hash starts with '$', invalid hash: "+hash)
			os.Exit(4)
		case len(hashpart[1]) == 0:
			fmt.Fprintln(os.Stderr, "Abort: Missing crypt type (between the 1st & 2nd '$'), invalid hash: "+hash)
			os.Exit(5)
		case hashpart[1][:1] != "2":
			fmt.Fprintln(os.Stderr, "Abort: The crypt type (between the 1st & 2nd '$') for bcrypt starts with '2', invalid hash: "+hash)
			os.Exit(6)
		case len(hashpart[1]) > 2:
			fmt.Fprintln(os.Stderr, "Abort: The crypt type (between the 1st & 2nd '$') is 1 or 2 characters long, invalid hash: "+hash)
			os.Exit(7)
		case len(hashpart[2]) != 2:
			fmt.Fprintln(os.Stderr, "Abort: The cost (between the 2nd & 3rd '$') must be 2 characters long, invalid hash: "+hash)
			os.Exit(8)
		case hashpart[2][0] < '0' || hashpart[2][0] > '3' || hashpart[2][1] < '0' || hashpart[2][1] > '9':
			fmt.Fprintln(os.Stderr, "Abort: The cost (between the 2nd & 3rd '$') must be numeric and 04..31, invalid hash: "+hash)
			os.Exit(9)
		case (hashpart[2][0] == '0' && hashpart[2][1] < '4') || (hashpart[2][0] == '3' && hashpart[2][1] > '1'):
			fmt.Fprintln(os.Stderr, "Abort: The cost (between the 2nd & 3rd '$') must be 04..31, invalid hash: "+hash)
			os.Exit(10)
		case len(hashpart[3]) != 53:
			fmt.Fprintln(os.Stderr, "Abort: The salt & password hash (after the 3rd '$') must be 53 characters long, invalid hash: "+hash)
			os.Exit(11)
		}
	}
	switch cmd {
	case "COST":
		if hash == "" {
			showHelp("Hash needed to tell its cost")
		}
		fmt.Printf("%d\n", getCost(hash))
	case "CHECK":
		doCheck(getPassword(), []byte(hash))
	case "HASH":
		if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
			showHelp("Argument for cost not 4..31, out of range: " + fmt.Sprint(cost))
		}
		fmt.Println(getHash(getPassword(), cost))
	}
}

func getCost(hash string) int {
	cost, e := bcrypt.Cost([]byte(hash))
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(12)
	}
	return cost
}

func doCheck(password, hash []byte) {
	if bcrypt.CompareHashAndPassword(hash, password) == nil {
		if !quiet {
			fmt.Println("true")
		}
	} else {
		if !quiet {
			fmt.Println("false")
		}
		os.Exit(1)
	}
}

func getHash(password []byte, cost int) string {
	hash, e := bcrypt.GenerateFromPassword(password, cost)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(13)
	}
	return string(hash)
}

func getPassword() []byte {
	var password []byte
	if !term.IsTerminal(0) {
		password, _ = io.ReadAll(os.Stdin)
		l := len(password)
		if l > 0 && password[l-1] == nl {
			password = password[:l-1]
		}
	} else {
		fmt.Fprintf(os.Stderr, "Enter password: ")
		pw, _ := term.ReadPassword(0)
		fmt.Fprintf(os.Stderr, "\r               \r")
		password = []byte(pw)
	}
	return password[:pwMaxLen]
}
