package main

import (
	"bufio"
	"os"
	"strings"
)

type db interface {
	set(name, val string) error
	get(name string) (val string, err error)
	delete(name string) (val string, err error) // TODO maybe just return err?
	count(val string) (count uint, err error)

	end() error

	begin() error
	rollback() error
	commit() error
}

type database struct {
	store map[string]string

	buf map[string]string
	txs []tx
}

type tx struct {
}

func err() {
	println("invalid input")
}

func main() {

	for {
		r := bufio.NewReader(os.Stdin)
		in, _ := r.ReadString('\n')

		args := strings.Split(in, " ")

		switch len(args) {
		case 1:
			switch strings.ToLower(args[0]) {
			case "begin\n":
			case "rollback\n":
			case "end\n": // TODO maybe strip trailing newline
				os.Exit(0)
			default:
				err()
			}
		case 2:
			switch strings.ToLower(args[0]) {
			case "get":
			case "delete":
			case "count":
			default:
				err()
			}
		case 3:
			if strings.ToLower(args[0]) == "set" {

			} else {
				err()
			}
		default:
			err()
		}
	}
}
