package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/atecce/devoted/db"
)

// TODO usage
func usage() {
	println("invalid input")
}

var dbg = flag.Bool("debug", false, "pretty print debug output")

func main() {

	flag.Parse()

	db := db.NewDatabase()

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
			db.Debug()
			println()
		}

		switch len(args) {
		case 1:
			switch strings.ToLower(args[0]) {
			case "begin":
				db.Begin()
			case "commit":
				db.Commit()
			case "rollback":
				db.Rollback()
			case "end":
				os.Exit(0)
			default:
				usage()
			}
		case 2:
			switch strings.ToLower(args[0]) {
			case "get":
				val := db.Get(args[1])
				if val == nil {
					fmt.Println(val)
				} else {
					fmt.Println(*val)
				}
			case "delete":
				db.Delete(args[1])
			case "count":
				fmt.Println(db.Count(args[1]))
			default:
				usage()
			}
		case 3:
			if strings.ToLower(args[0]) == "set" {
				db.Set(args[1], args[2])
			} else {
				usage()
			}
		default:
			usage()
		}

		if *dbg {
			println()
			println("after...")
			db.Debug()
			println()
		}
	}
}
