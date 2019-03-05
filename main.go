package main

import (
	"bufio"
	"errors"
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

func (db database) get(name string) (*string, error) {
	if val, ok := db.store[name]; ok {
		return &val, nil
	}
	return nil, errors.New("not found")
}

func (db database) set(name, val string) error {
	db.store[name] = val
	return nil
}

type tx struct {
}

// TODO usage
func err() {
	println("invalid input")
}

func main() {

	db := database{}
	db.store = make(map[string]string)

	for {
		r := bufio.NewReader(os.Stdin)
		in, _ := r.ReadString('\n')

		stripped := strings.TrimSuffix(in, "\n")

		args := strings.Split(stripped, " ")

		switch len(args) {
		case 1:
			switch strings.ToLower(args[0]) {
			case "begin":
			case "rollback":
			case "end": // TODO maybe strip trailing newline
				os.Exit(0)
			default:
				err()
			}
		case 2:
			switch strings.ToLower(args[0]) {
			case "get":
				// pretty.Println(db.store)
				val, err := db.get(args[1])
				if err != nil {
					println("<nil>") // TODO this could actually be a val
				} else {
					os.Stdout.Write([]byte(*val + "\n")) // TODO could be nil?
				}
			case "delete":
			case "count":
			default:
				err()
			}
		case 3:
			if strings.ToLower(args[0]) == "set" {
				if err := db.set(args[1], args[2]); err != nil {
					println("this shouldn't happen")
				}
				// pretty.Println(db.store)
			} else {
				err()
			}
		default:
			err()
		}
	}
}
