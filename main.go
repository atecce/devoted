package main

import (
	"bufio"
	"os"
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
		txt, _ := r.ReadString('\n')
		print(txt)
	}
}
