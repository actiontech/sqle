package go_ora

type Transaction struct {
	conn *Connection
}

func (tx *Transaction) Commit() error {
	tx.conn.autoCommit = true
	tx.conn.session.ResetBuffer()
	return (&simpleObject{connection: tx.conn, operationID: 0xE}).write().read()
}

func (tx *Transaction) Rollback() error {
	tx.conn.autoCommit = true
	tx.conn.session.ResetBuffer()
	return (&simpleObject{connection: tx.conn, operationID: 0xF}).write().read()
}
