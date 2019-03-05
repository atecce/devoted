package main

type database struct {
	store map[string]string

	txs *[]transaction
}

type transaction struct {
	setBuf    map[string]string
	deleteBuf *[]string
}

func newDatabase() database {
	return database{
		store: make(map[string]string),
		txs:   new([]transaction),
	}
}

func (db *database) get(name string) *string {

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

func (db *database) set(name, val string) {
	n := len(*db.txs)
	if n != 0 {
		tx := (*db.txs)[n-1]
		tx.setBuf[name] = val
	} else {
		db.store[name] = val
	}
}

func (db *database) delete(name string) {

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

func (db *database) count(val string) uint {

	var n uint

	// go backwards in time
	var futureDeletes []string
	for i := len(*db.txs) - 1; i >= 0; i-- {

		tx := (*db.txs)[i]

		// remember future deletes
		futureDeletes = append(futureDeletes, *tx.deleteBuf...)

		for k, v := range tx.setBuf {
			if v == val {

				// if we are going to delete this name in the future
				// don't count it
				if contains(futureDeletes, k) {
					continue
				}

				n++
			}
		}
	}

	for k, v := range db.store {
		if val == v {

			// if we are going to delete this name in the future
			// don't count it
			if contains(futureDeletes, k) {
				continue
			}

			n++
		}
	}

	return n
}

func (db *database) begin() {
	*db.txs = append(*db.txs, transaction{
		setBuf:    make(map[string]string),
		deleteBuf: new([]string),
	})
}

func (db *database) rollback() {
	if n := len(*db.txs); n > 0 {
		*db.txs = (*db.txs)[:n-1]
	} else {
		// no transactions
	}
}

func (db *database) commit() {

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

		// reset transactions
		db.txs = new([]transaction)
	}
}
