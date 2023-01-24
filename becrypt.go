// becrypt - CLI tool for generating and checking bcrypt hashes
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	//"golang.org/x/crypto/ssh/terminal"
)

const (
	version = "1.2.1"
	passwordMaxLen = 72
)

var (
	helptxt = ""
	usage = ` - CLI tool for generating and checking bcrypt hashes
Usage:
becrypt [<cost>] | <hash> | cost <hash>
    becrypt [<cost>]:     Generate a hash from the password
                          (optional <cost>: ` + strconv.Itoa(bcrypt.MinCost) +
		".." + strconv.Itoa(bcrypt.MaxCost) + ", default: " +
		strconv.Itoa(bcrypt.DefaultCost) + `)
    becrypt <hash>:       Check the password against <hash>
    becrypt cost <hash>:  Display the cost of <hash>
    becrypt help:         Display this help text
  Passwords can be piped-in or prompted for, maximum length: ` +
		strconv.Itoa(passwordMaxLen) + " characters."
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
		if arg == "cost" {
			cmd = arg
			continue
		}
		if arg == "help" {
			showHelp("")
		}
		if cmd == "" {
			cmd = "hash"
			n, _ := fmt.Sscan(arg, &cost)
			if n < 1 { // Not a hash but a check command
				cmd, hash = "check", arg
			}
			continue
		}
		showHelp("Too many arguments for " + cmd)
	}
	switch cmd {
	case "cost":
		fmt.Println(fmt.Sprintf("%d", getCost(hash)))
	case "check":
		if check(getPassword(), hash) {
			fmt.Println("yes")
		} else {
			fmt.Println("no")
			os.Exit(1)
		}
	case "hash":
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
	c, e := bcrypt.Cost([]byte(hash))
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(3)
	}
	return c
}

func check(password, hash string) bool {
	e := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return false
	}
	return true
}

func getHash(password string, cost int) string {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		showHelp("Argument for cost out of range")
	}
	hash, e := bcrypt.GenerateFromPassword([]byte(getPassword()), cost)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(4)
	}
	return string(hash)
}

func getPassword() string {
	password, e := io.ReadAll(os.Stdin)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(5)
	}
	if password == nil {
		retryTimes := 3
		for retryTimes > 0 {
			fmt.Printf("Enter password: ")
			pwd, _ := term.ReadPassword(0)
			if len(pwd) == 0 {
				fmt.Printf("\nPassword can't be empty")
			} else {
				fmt.Printf("\nConfirm password: ")
				pwdc, _ := term.ReadPassword(0)
				fmt.Println()
				if bytes.Equal(pwd, pwdc) {
					password = pwd
					break
				} else {
					fmt.Printf("Passwords not the same")
				}
			}
			retryTimes--
			if retryTimes > 0 {
				fmt.Printf(", retry")
			}
		}
		if password == nil {
			showHelp("Password missing")
		}
	}
	if len([]byte(password)) > passwordMaxLen {
		showHelp("Password cannot be longer than " + strconv.Itoa(passwordMaxLen) + " bytes")
	}
	return string(password)
}
