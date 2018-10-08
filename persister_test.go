package otira

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestPersistSimple(t *testing.T) {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	p, err := NewPersister(db, new(DialectSqlite3), 10)
	if err != nil {
		t.Fatal(err)
	}
	close(p.incoming)
}

func TestPersistFewRecords(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	p, err := NewPersister(db, new(DialectSqlite3), 10)
	if err != nil {
		t.Fatal(err)
	}

	table, err := defaultTestTable()
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		tableRecord, err := table.NewRecord()
		if err != nil {
			t.Fatal(err)
		}
		p.incoming <- tableRecord
	}

	close(p.incoming)
}

///// FAILS /////
func TestPersistNoDbFail(t *testing.T) {
	_, err := NewPersister(nil, new(DialectSqlite3), 10)
	if err == nil {
		t.Fatal(err)
	}
}

func TestNoDialectFail(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = NewPersister(db, nil, 10)
	if err == nil {
		t.Fatal(err)
	}
}

func TestPersistBadNumFail(t *testing.T) {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = NewPersister(db, new(DialectSqlite3), -1)
	if err == nil {
		t.Fatal(err)
	}

}
