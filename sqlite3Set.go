package otira

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type sqlite3Set struct {
	dir     string
	counter uint32
	pers    *Persister
	table   *TableDef
}

func NewSqlite3Set(dir string) (Set, error) {
	if dir == "" {
		return nil, errors.New("HashCache dir is empty")
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		err = os.RemoveAll(dir)
		if err != nil {
			return nil, err
		}
	}
	err := os.Mkdir(dir, 0755)
	if err != nil {
		return nil, err
	}

	set := new(sqlite3Set)
	set.dir = dir

	db, err := sql.Open("sqlite3", dir+"/foom6")
	if err != nil {
		return nil, err
	}
	set.pers, err = NewPersister(db, new(DialectSqlite3), 1000)
	if err != nil {
		return nil, err
	}

	set.table, err = NewTableDef("set")
	if err != nil {
		return nil, err
	}
	id := new(FieldDefInt64)
	id.SetName("id")
	id.SetUnique(true)
	err = set.table.Add(id)
	if err != nil {
		return nil, err
	}
	set.table.Done()
	set.pers.CreateTables(set.table)

	set.counter = 0
	if err != nil {
		return nil, err
	}
	return set, nil

}

func (set *sqlite3Set) Close() error {
	err := os.RemoveAll(set.dir)
	if err != nil {
		return err
	}
	return nil
}

func (set *sqlite3Set) Contains(key int64) (bool, error) {
	s := "select count(id) from set where id=?"
	stmt, err := set.pers.tx.Prepare(s)
	defer stmt.Close()
	if err != nil {
		return false, err
	}

	var count int64
	err = stmt.QueryRow(key).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (set *sqlite3Set) Put(key int64) error {
	set.counter++
	rec, err := set.table.NewRecord()
	rec.values[0] = key
	err = set.pers.save(rec)
	if err != nil {
		return err
	}
	return nil
}
