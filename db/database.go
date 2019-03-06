package db

import "github.com/kr/pretty"

type Database struct {
	store map[string]string

	txs *[]transaction
}

type transaction struct {
	setBuf    map[string]string
	deleteBuf *[]string
}

func NewDatabase() Database {
	return Database{
		store: make(map[string]string),
		txs:   new([]transaction),
	}
}

func (db *Database) Get(name string) *string {

	// go backwards in time
	for i := len(*db.txs) - 1; i >= 0; i-- {
		tx := (*db.txs)[i]

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

	// check the "persisted" store
	if val, ok := db.store[name]; ok {
		return &val
	}

	return nil
}

func (db *Database) Set(name, val string) {
	if n := len(*db.txs); n == 0 {
		// if there are no transactions, set in the store
		db.store[name] = val
	} else {
		// otherwise set in the most recent buffer
		tx := (*db.txs)[n-1]
		tx.setBuf[name] = val
	}
}

func (db *Database) Delete(name string) {

	if n := len(*db.txs); n != 0 {
		tx := (*db.txs)[n-1]

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

func (db *Database) Count(val string) uint {

	var n uint

	// go backwards in time
	var futureDeletes []string
	for i := len(*db.txs) - 1; i >= 0; i-- {

		tx := (*db.txs)[i]

		// remember future deletes
		futureDeletes = append(futureDeletes, *tx.deleteBuf...)

		n += count(val, tx.setBuf, futureDeletes)
	}

	n += count(val, db.store, futureDeletes)

	return n
}

func count(val string, m map[string]string, skip []string) uint {
	var n uint

	for k, v := range m {
		if val == v {

			// if we are going to delete this name in the future
			// don't count it
			if contains(skip, k) {
				continue
			}

			n++
		}
	}
	return n
}

func contains(strs []string, val string) bool {
	for _, str := range strs {
		if str == val {
			return true
		}
	}
	return false
}

func (db *Database) Begin() {
	*db.txs = append(*db.txs, transaction{
		setBuf:    make(map[string]string),
		deleteBuf: new([]string),
	})
}

func (db *Database) Rollback() {
	if n := len(*db.txs); n == 0 {
		println("no transactions to rollback")
	} else {
		*db.txs = (*db.txs)[:n-1]
	}
}

func (db *Database) Commit() {

	// go backwards in time
	for i := len(*db.txs) - 1; i >= 0; i-- {
		tx := (*db.txs)[i]

		// persist buffers to store
		for k, v := range tx.setBuf {
			db.store[k] = v
		}
		for _, k := range *tx.deleteBuf {
			delete(db.store, k)
		}
	}

	// reset transactions
	db.txs = new([]transaction)
}

func (db *Database) Debug() {
	pretty.Println("store:", db.store)
	pretty.Println("txs:", db.txs)
}
