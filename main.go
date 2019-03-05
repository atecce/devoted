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

func main() {

	for {
		r := bufio.NewReader(os.Stdin)
		in, _ := r.ReadString('\n')

		args := strings.Split(in, " ")

		switch len(args) {
		case 1:
			switch strings.ToLower(args[0]) {
			case "end\n": // TODO maybe strip trailing newline
				os.Exit(0)
			default:
				println("invalid input")
			}
		case 2:
			print(args[0], args[1])
		case 3:
			print(args[0], args[1], args[2])
		default:
			println("invalid input")
		}
	}
}
