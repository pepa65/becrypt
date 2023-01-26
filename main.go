// becrypt - CLI tool for generating and checking bcrypt hashes
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
	version = "1.2.2"
	pwMaxLen = 72
)

var (
	helptxt = ""
	usage = ` - CLI tool for generating and checking bcrypt hashes
Repo:   github.com/pepa65/becrypt
Usage:  becrypt [<cost>] | <hash> | cost|-c|--cost <hash> | help|-h|--help
    becrypt [<cost>]:               Generate a hash from the password
                                    (optional <cost>: ` +
		strconv.Itoa(bcrypt.MinCost) + ".." + strconv.Itoa(bcrypt.MaxCost) +
		", default: " + strconv.Itoa(bcrypt.DefaultCost) + `)
    becrypt <hash>:                 Check the password against <hash>
    becrypt cost|-c|--cost <hash>:  Display the cost of <hash>
    becrypt help|-h|--help:         Display this help text
  The password can be piped-in or prompted for, is cut off after ` +
		strconv.Itoa(pwMaxLen) + " characters."
)

func main() { // I:self,version,usage
	self, hash, cmd, cost := "", "", "", bcrypt.DefaultCost
	for _, arg := range os.Args {
		if self == "" {
			selves := strings.Split(arg, "/")
			self = selves[len(selves)-1]
			helptxt = self + " v" + version + usage
			continue
		}
		if cmd == "cost" {
			if hash != "" {
				showHelp("Too many arguments for cost")
			}
			hash = arg
			continue
		}
		if cmd == "" && (arg == "cost" || arg == "-c" || arg == "--cost") {
			cmd = "cost"
			continue
		}
		if arg == "help" || arg == "-h" || arg == "--help" {
			showHelp("")
		}
		if cmd == "" {
			c, e := strconv.Atoi(arg)
			if e == nil { // Integer: a hash command
				cmd, cost = "hash", c
				continue
			}
			cmd, hash = "check", arg
		} else {
			showHelp("Too many arguments")
		}
	}
	if cmd == "" {
		cmd = "hash"
	}
	//fmt.Fprintln(os.Stderr, "Command: " + cmd)
	switch cmd {
	case "cost":
		fmt.Printf("%d\n", getCost(hash))
	case "check":
		doCheck(getPassword(), []byte(hash))
	case "hash":
		if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
			showHelp("Argument for cost out of range")
		}
		fmt.Println(getHash(getPassword(), cost))
	}
}

func showHelp(msg string) { // I:helptxt
	fmt.Fprintln(os.Stderr, helptxt)
	if msg != "" {
		fmt.Fprintln(os.Stderr, "Abort: " + msg)
		os.Exit(2)
	}
	os.Exit(0)
}

func getCost(hash string) int {
	cost, e := bcrypt.Cost([]byte(hash))
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(3)
	}
	return cost
}

func doCheck(password, hash []byte) {
	if bcrypt.CompareHashAndPassword(hash, password) == nil {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
		os.Exit(1)
	}
}

func getHash(password []byte, cost int) string {
	hash, e := bcrypt.GenerateFromPassword(password, cost)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(4)
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
