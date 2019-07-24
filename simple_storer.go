package otira

import (
	"database/sql"
	"errors"
)

type SimpleStorer struct {
	db        *sql.DB
	dbopts    *DBOptions
	tables    []*TableDef
	persister *Persister
}

func (st *SimpleStorer) Initialize(db *sql.DB, dbopts *DBOptions, tms ...*TableDef) error {
	if db == nil || tms == nil {
		return errors.New("Cannot have nil db or TableDef")
	}

	pers, err := NewPersister(db, dbopts.Dialect, dbopts.TxSize)
	if err != nil {
		return err
	}

	st.dbopts = dbopts

	if st.dbopts != nil && st.dbopts.CreateTables {

		err = pers.CreateTables(tms...)
		if err != nil {
			return err
		}
	}
	err = pers.BeginTx()
	if err != nil {
		return err
	}

	st.db = db
	st.dbopts = dbopts
	st.tables = tms
	st.persister = pers

	return err
}

func (st *SimpleStorer) Save(rec *Record) error {
	err := st.persister.Save(rec)
	return err
}

func (st *SimpleStorer) Close() error {
	return st.persister.Done()
}
