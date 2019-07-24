package otira

import (
	"database/sql"
)

type DBOptions struct {
	Opts         *sql.TxOptions
	Dialect      Dialect
	TxSize       int
	ChannelSize  int
	CreateTables bool
}

type Storer interface {
	Initialize(db *sql.DB, dbopts *DBOptions, tms ...*TableDef) error
	Save(*Record) error
	Close() error
}
