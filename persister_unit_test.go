package otira

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestNewPersist_NilDb(t *testing.T) {
	_, err := NewPersister(nil, new(DialectSqlite3))
	if err == nil {
		t.Fatal(err)
	}
}

func TestNewPersist_NilDialect(t *testing.T) {
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

func TestNewPersist_CreateTablesNilList(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	pers, err := NewPersister(db, new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(nil)
	if err == nil {
		t.Fatal(err)
	}

}

func TestNewPersist_CreateTableNil(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	pers, err := NewPersister(db, new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}

	err = pers.createTable(nil)
	if err == nil {
		t.Fatal(err)
	}

}
