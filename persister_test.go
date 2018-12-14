package otira

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestPersistInstantiate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	pers, err := NewPersister(db, new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}
	pers.Done()

}

func TestPersistFewRecords(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	//db, err := sql.Open("sqlite3", "/tmp/mmmmmm9")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	pers, err := NewPersister(db, new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}
	defer pers.Done()

	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Error(err)
	}

	table.SetDone()

	err = pers.CreateTables(table)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 100; i++ {
		//t.Log(i)
		rec, err := table.NewRecord()
		if err != nil {
			t.Fatal(err)
		}
		populateDefaultTableRecord(rec)
		err = rec.SetByName(pk, 1000000+i)
		if err != nil {
			t.Fatal(err)
		}
		err = pers.save(rec)
		if err != nil {
			t.Fatal(err)
		}
	}

}

///// FAILS /////
func TestPersistNoDbFail(t *testing.T) {
	_, err := NewPersister(nil, new(DialectSqlite3))
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

	_, err = NewPersister(db, nil)
	if err == nil {
		t.Fatal(err)
	}
}

// func TestPersistBadNumFail(t *testing.T) {
// 	db, err := sql.Open("sqlite3", ":memory:")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()

// 	_, err = NewPersister(db, new(DialectSqlite3))
// 	if err == nil {
// 		t.Fatal(err)
// 	}
// }
