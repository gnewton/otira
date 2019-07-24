package otira

import (
	"context"
	"database/sql"
	"errors"
	"sync"
)

type TransactionHandler struct {
	ctx        *context.Context
	dialect    *Dialect
	conn       *sql.Conn
	opts       *sql.TxOptions
	tx         *sql.Tx
	counter    int
	txSize     int
	writeMutex sync.Mutex
	c          chan []*Record
}

func NewTransactionHandler(c chan []*Record, conn *sql.Conn, ctx *context.Context, opts *sql.TxOptions, dialect *Dialect, txSize int) (*TransactionHandler, error) {
	if conn == nil {
		return nil, errors.New("conn is nil")
	}

	th := TransactionHandler{
		conn:    conn,
		txSize:  txSize,
		dialect: dialect,
		ctx:     ctx,
		opts:    opts,
		c:       c,
	}

	error := th.reset()

	return &th, error
}

func (th *TransactionHandler) commit() error {
	err := th.tx.Commit()
	return err
}

func (th *TransactionHandler) start() {

}

func (th *TransactionHandler) reset() error {
	th.counter = 0
	var err error
	th.tx, err = th.conn.BeginTx(*th.ctx, th.opts)
	return err

}

func (th *TransactionHandler) save(records []Record) error {
	th.writeMutex.Lock()
	defer th.writeMutex.Unlock()

	if records == nil {
		return errors.New("Records is nil")
	}
	recLen := len(records)
	if recLen == 0 {
		return errors.New("Empty records")
	}

	th.counter += recLen
	if recLen > th.txSize {
		err := th.commit()
		if err != nil {
			return err
		}
		err = th.reset()
		if err != nil {
			return err
		}

	}

	return nil
}
