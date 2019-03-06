package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kr/pretty"
)

// TODO usage
func usage() {
	println("invalid input")
}

func debug(db database) {
	pretty.Println("store:", db.store)
	pretty.Println("txs:", db.txs)
}

var dbg = flag.Bool("debug", false, "pretty print debug output")

func main() {

	flag.Parse()

	db := newDatabase()

	r := bufio.NewReader(os.Stdin)

	for {

		print(">> ")
		in, err := r.ReadString('\n')
		if err != nil {
			println("failed to read string. quitting")
			os.Exit(1)
		}

		args := strings.Split(strings.TrimSuffix(in, "\n"), " ")

		if *dbg {
			println()
			println("before...")
			debug(db)
			println()
		}

		switch len(args) {
		case 1:
			switch strings.ToLower(args[0]) {
			case "begin":
				db.begin()
			case "commit":
				db.commit()
			case "rollback":
				db.rollback()
			case "end":
				os.Exit(0)
			default:
				usage()
			}
		case 2:
			switch strings.ToLower(args[0]) {
			case "get":
				val := db.get(args[1])
				if val == nil {
					fmt.Println(val)
				} else {
					fmt.Println(*val)
				}
			case "delete":
				db.delete(args[1])
			case "count":
				fmt.Println(db.count(args[1]))
			default:
				usage()
			}
		case 3:
			if strings.ToLower(args[0]) == "set" {
				db.set(args[1], args[2])
			} else {
				usage()
			}
		default:
			usage()
		}

		if *dbg {
			println()
			println("after...")
			debug(db)
			println()
		}
	}
}
