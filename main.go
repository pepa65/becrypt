// becrypt - Generate and test bcrypt hashes from a CLI

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
	version = "1.2.3"
	pwMaxLen = 72
)

var (
	self = ""
)

func showHelp(msg string) { // I:self,version,pwMaxLen
	helptxt := self + " v" + version + ` - Generate and test bcrypt hashes from a CLI
Repo:   github.com/pepa65/becrypt
Usage:  ` + self + ` HASH | TEST | COST | HELP
 HASH:  [<cost>]:               Generate a hash from the password
                               (optional <cost>: ` +
    strconv.Itoa(bcrypt.MinCost) + ".." + strconv.Itoa(bcrypt.MaxCost) +
    ", default: " + strconv.Itoa(bcrypt.DefaultCost) + `)
 TEST:  <hash>:                 Test the password against <hash>
 COST:  cost|-c|--cost <hash>:  Display the cost of <hash>
 HELP:  help|-h|--help:         Display this help text
The password can be piped-in or prompted for, is cut off after ` +
    strconv.Itoa(pwMaxLen) + " characters."
	fmt.Fprintln(os.Stderr, helptxt)
	if msg != "" {
		fmt.Fprintln(os.Stderr, "Abort: " + msg)
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
				showHelp("Too many arguments for cost")
			}
			hash = arg
			continue
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
			cmd, hash = "TEST", arg
		} else {
			showHelp("Too many arguments")
		}
	}
	if cmd == "" {
		cmd = "HASH"
	}
	if hash != "" { // test hash
		hashpart := strings.Split(hash, "$")
		switch true {
		case len(hashpart) != 4:
			fmt.Fprintln(os.Stderr, "Abort: Exactly 3 '$' in a bcrypt hash")
			os.Exit(10)
		case len(hashpart[0]) > 0:
			fmt.Fprintln(os.Stderr, "Abort: A proper bcrypt hash starts with '$'")
			os.Exit(3)
		case len(hashpart[1]) == 0:
			fmt.Fprintln(os.Stderr, "Abort: Missing crypt type (between the 1st & 2nd '$')")
			os.Exit(4)
		case hashpart[1][:1] != "2":
			fmt.Fprintln(os.Stderr, "Abort: The crypt type (between the 1st & 2nd '$') for bcrypt starts with '2'")
			os.Exit(5)
		case len(hashpart[1]) > 2:
			fmt.Fprintln(os.Stderr, "Abort: The crypt type (between the 1st & 2nd '$') is 1 or 2 characters long")
			os.Exit(6)
		case len(hashpart[2]) != 2:
			fmt.Fprintln(os.Stderr, "Abort: The cost (between the 2nd & 3rd '$') must be 2 characters long")
			os.Exit(7)
		case hashpart[2][0] < '0' || hashpart[2][0] > '3' || hashpart[2][1] < '0' || hashpart[2][1] > '9':
			fmt.Fprintln(os.Stderr, "Abort: The cost (between the 2nd & 3rd '$') must be numeric and 04..31")
			os.Exit(8)
		case (hashpart[2][0] == '0' && hashpart[2][1] < '4') || (hashpart[2][0] == '3' && hashpart[2][1] > '1'):
			fmt.Fprintln(os.Stderr, "Abort: The cost (between the 2nd & 3rd '$') must be 04..31")
			os.Exit(9)
		case len(hashpart[3]) != 53:
			fmt.Fprintln(os.Stderr, "Abort: The salt & password hash (after the 3rd '$') must be 53 characters long")
			os.Exit(10)
		}
	}
	switch cmd {
	case "COST":
		fmt.Printf("%d\n", getCost(hash))
	case "TEST":
		doTest(getPassword(), []byte(hash))
	case "HASH":
		if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
			showHelp("Argument for cost out of range")
		}
		fmt.Println(getHash(getPassword(), cost))
	}
}

func getCost(hash string) int {
	cost, e := bcrypt.Cost([]byte(hash))
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(11)
	}
	return cost
}

func doTest(password, hash []byte) {
	if bcrypt.CompareHashAndPassword(hash, password) == nil {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
		os.Exit(12)
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
	if ! term.IsTerminal(0) {
		password, _ = io.ReadAll(os.Stdin)
	}
	if len(password) == 0 {
		fmt.Fprintf(os.Stderr, "Enter password: ")
		pw, _ := term.ReadPassword(0)
		fmt.Fprintf(os.Stderr, "\r               \r")
		password = []byte(pw)
	}
	if len(password) > 72 {
		password = password[:pwMaxLen]
	}
	return password
}
