package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kr/pretty"
)

type database struct {
	store map[string]string

	tx        bool
	setBuf    map[string]string
	deleteBuf []string
}

func (db database) get(name string) *string {
	// TODO behavior during a tx?
	if val, ok := db.store[name]; ok {
		return &val
	}
	return nil
}

func (db database) set(name, val string) {
	if db.tx {
		db.setBuf[name] = val
	} else {
		db.store[name] = val
	}
}

func (db database) delete(name string) {
	if db.tx {
		db.deleteBuf = append(db.deleteBuf, name)
	} else {
		delete(db.store, name)
	}
}

func (db database) count(val string) uint {
	var count uint
	for _, value := range db.store {
		if val == value {
			count++
		}
	}
	return count
}

func (db database) rollback() error {
	db.setBuf = make(map[string]string)
	db.deleteBuf = nil
	db.tx = false
	return nil
}

func (db database) commit() error {
	for name, val := range db.setBuf {
		db.store[name] = val
	}
	for _, name := range db.deleteBuf {
		// TODO maybe drain buffer first
		delete(db.store, name)
	}
	db.setBuf = make(map[string]string)
	db.deleteBuf = nil
	return nil
}

// TODO usage
func err() {
	println("invalid input")
}

func main() {

	db := database{}
	db.store = make(map[string]string)
	db.setBuf = make(map[string]string)

	r := bufio.NewReader(os.Stdin)

	for {
		in, _ := r.ReadString('\n')

		stripped := strings.TrimSuffix(in, "\n")

		args := strings.Split(stripped, " ")

		switch len(args) {
		case 1:
			switch strings.ToLower(args[0]) {
			case "begin":
				db.tx = true
			case "commit":
				db.commit()
				db.tx = false
			case "rollback":
				db.rollback()
				db.tx = false
			case "end":
				os.Exit(0)
			default:
				err()
			}
		case 2:
			switch strings.ToLower(args[0]) {
			case "get":
				pretty.Println(db.store)
				val := db.get(args[1])
				if val == nil {
					println("<nil>") // TODO this could actually be a val
				} else {
					fmt.Println(*val)
				}
			case "delete":
				db.delete(args[1])
			case "count":
				fmt.Println(db.count(args[1]))
			default:
				err()
			}
		case 3:
			if strings.ToLower(args[0]) == "set" {
				db.set(args[1], args[2])
				pretty.Println(db.store)
			} else {
				err()
			}
		default:
			err()
		}
	}
}
