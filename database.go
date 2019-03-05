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
