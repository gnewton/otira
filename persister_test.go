package otira

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestPersistSimple(t *testing.T) {
	t.Log("hello")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	pers, err := NewPersister(db, new(DialectSqlite3), 10)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("hello")
	pers.Done()
	t.Log("fpp hello")

}

func TestPersistFewRecords(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	pers, err := NewPersister(db, new(DialectSqlite3), 10)
	if err != nil {
		t.Fatal(err)
	}
	defer pers.Done()

	table, err := defaultTestTable()
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 100; i++ {
		tableRecord, err := table.NewRecord()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("TestPersistFewRecords")
		t.Log(*tableRecord)

		err = pers.Save(tableRecord)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestPersistFewRecordsWithCancel(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	pers, err := NewPersister(db, new(DialectSqlite3), 10)
	if err != nil {
		t.Fatal(err)
	}

	table, err := defaultTestTable()
	if err != nil {
		t.Error(err)
	}
	table.name = "foobar"
	for i := 0; i < 100; i++ {
		tableRecord, err := table.NewRecord()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("TestPersistFewRecordsWithCancel")
		t.Log(*tableRecord)
		pers.Save(tableRecord)
		if i == 5 {
			pers.cancelFunc()
			break
		}
	}
	pers.Done()
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
