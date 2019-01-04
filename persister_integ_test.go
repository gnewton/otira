package otira

import (
	"database/sql"
	"errors"
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
		t.Fatal(err)
	}

	table.SetDone()

	err = pers.CreateTables(table)
	if err != nil {
		t.Fatal(err)
	}

	var i uint64
	base := uint64(100000)
	for i = 0; i < 100; i++ {
		//t.Log(i)
		rec, err := table.NewRecord()
		if err != nil {
			t.Fatal(err)
		}
		populateDefaultTableRecord(rec)
		err = rec.SetByName(pk, base+i)
		//err = rec.SetByName(pk, "foo")
		if err != nil {
			t.Fatal(err)
		}
		err = pers.save(rec)
		if err != nil {
			t.Fatal(err)
		}
		exists, err := pers.Exists(rec)
		if err != nil {
			t.Fatal(err)
		}
		if !exists {
			t.Fatal(errors.New("Should exist"))
		}
	}

}

func TestSaveThenUpdateRecord(t *testing.T) {
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
		t.Fatal(err)
	}

	table.SetDone()

	err = pers.CreateTables(table)
	if err != nil {
		t.Fatal(err)
	}
	base := uint64(100000)

	rec, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}
	populateDefaultTableRecord(rec)
	err = rec.SetByName(pk, base)
	if err != nil {
		t.Fatal(err)
	}
	err = pers.save(rec)
	if err != nil {
		t.Fatal(err)
	}
	exists, err := pers.Exists(rec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Update")
	err = rec.SetByName(tAddress, "UPDATE street test")
	if err != nil {
		t.Fatal(err)
	}

	//err = pers.save(rec)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal(errors.New("Should exist"))
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

func TestPersistPragmas(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	//db, err := sql.Open("sqlite3", "foo.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	pers, err := NewPersister(db, new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	//	err = pers.initPragmas()
	if err != nil {
		t.Fatal(err)
	}

	pers.Done()

}
