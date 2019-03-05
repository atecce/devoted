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

	txs *[]transaction
}

type transaction struct {
	setBuf    map[string]string
	deleteBuf *[]string
}

func (db *database) get(name string) *string {

	txsLen := len(*db.txs)
	if txsLen != 0 {

		tx := (*db.txs)[txsLen-1]

		// check if we deleted it in this tx
		for _, delName := range *tx.deleteBuf {
			if name == delName {
				return nil
			}
		}

		// check if we set it in this tx
		if val, ok := tx.setBuf[name]; ok {
			return &val
		}

	}

	if val, ok := db.store[name]; ok {
		return &val
	}

	return nil
}

func (db *database) set(name, val string) {

	txsLen := len(*db.txs)
	if txsLen != 0 {
		tx := (*db.txs)[txsLen-1]
		tx.setBuf[name] = val
	} else {
		db.store[name] = val
	}
}

func (db *database) delete(name string) {

	// TODO no-op proper semantics when name doesn't exist?

	txsLen := len(*db.txs)
	if txsLen != 0 {
		tx := (*db.txs)[txsLen-1]

		// check if we set the key in this transaction
		if _, ok := tx.setBuf[name]; ok {
			delete(tx.setBuf, name)
			return
		}

		// otherwise queue up a delete
		*tx.deleteBuf = append(*tx.deleteBuf, name)

	} else {
		delete(db.store, name)
	}
}

func (db *database) count(val string) uint {

	var count uint

	// go backwards in time
	var futureDeletedNames []string
	for i := len(*db.txs) - 1; i >= 0; i-- {

		tx := (*db.txs)[i]

		// save future deleted names
		futureDeletedNames = append(futureDeletedNames, *tx.deleteBuf...)

		for name, value := range tx.setBuf {
			if value == val {
				// if we are going to delete this name in the future
				// don't count it
				if contains(futureDeletedNames, name) {
					continue
				}

				count++
			}
		}
	}

	for name, value := range db.store {
		if val == value {
			// if we are going to delete this name in the future
			// don't count it
			if contains(futureDeletedNames, name) {
				continue
			}

			count++
		}
	}

	return count
}

func (db *database) begin() {
	*db.txs = append(*db.txs, transaction{
		setBuf:    make(map[string]string),
		deleteBuf: new([]string),
	})
}

func (db *database) rollback() {
	txsLen := len(*db.txs)
	if txsLen > 0 {
		*db.txs = (*db.txs)[:txsLen-1]
	} else {
		// no transactions
	}
}

func (db *database) commit() {

	txsLen := len(*db.txs)
	if txsLen != 0 {
		tx := (*db.txs)[txsLen-1]

		for name, val := range tx.setBuf {
			db.store[name] = val
		}
		for _, name := range *tx.deleteBuf {
			// TODO maybe drain buffer first
			delete(db.store, name)
		}
	} else {
		// no transactions
	}
}

// TODO usage
func err() {
	println("invalid input")
}

func debug(db database) {
	pretty.Println("store:", db.store)
	pretty.Println("txs:", db.txs)
}

func main() {

	db := database{
		store: make(map[string]string),
		txs:   new([]transaction),
	}

	r := bufio.NewReader(os.Stdin)

	for {
		in, _ := r.ReadString('\n')

		stripped := strings.TrimSuffix(in, "\n")

		args := strings.Split(stripped, " ")

		// TODO add debug flag
		// println()
		// println("before...")
		// debug(db)

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
				err()
			}
		case 2:
			switch strings.ToLower(args[0]) {
			case "get":
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
			} else {
				err()
			}
		default:
			err()
		}

		// TODO add debug flag
		// println()
		// println("after...")
		// debug(db)

		// println()
	}
}

func contains(strs []string, val string) bool {
	for _, str := range strs {
		if str == val {
			return true
		}
	}
	return false
}
