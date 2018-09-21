package otira

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreatePrepared(t *testing.T) {
	table, err := defaultTestTable()
	if err != nil {
		t.Error(err)
	}

	prep, err := table.CreatePreparedStatementInsertAllFields(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(prep)
}

func TestTableCreate(t *testing.T) {
	table, err := defaultTestTable()
	if err != nil {
		t.Error(err)
	}

	cr, err := table.CreateTableString(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cr)

}

func TestRecordFromTableMeta(t *testing.T) {
	table, err := defaultTestTable()
	if err != nil {
		t.Error(err)
	}
	record, err := table.NewRecordSomeFields(table.fields[0], table.fields[1], table.fields[2])
	record.validating = true

	if err != nil {
		t.Error(err)
	}
	err = record.Set(0, 42)
	if err != nil {
		t.Error(err)
	}

	err = record.Set(1, "mm")
	if err != nil {
		t.Error(err)
	}

	err = record.Set(2, 4.5)
	if err != nil {
		t.Error(err)
	}
	prepared, err := table.CreatePreparedStatementInsertFromRecord(new(DialectSqlite3), record)
	log.Println("Prepared", prepared)

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE journals (id INTEGER PRIMARY KEY, firstname CHAR(24), age REAL, address VARCHAR(64), height REAL)")

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := tx.Prepare(prepared)
	if err != nil {
		t.Fatal(err)
	}

	result, err := stmt.Exec(record.values...)
	record.Set(0, 55)
	result, err = stmt.Exec(record.values...)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.RowsAffected())
	tx.Commit()

}

func TestRecordFromTableMetaTODO(t *testing.T) {
	table, err := defaultTestTable()
	if err != nil {
		t.Error(err)
	}
	record, err := table.NewRecord()

	err = record.Set(0, "mm")
	if err != nil {
		t.Error(err)
	}
	v := 44
	err = record.Set(1, v)
	if err != nil {
		t.Error(err)
	}
	log.Printf("\nRecord: %v", *record)
	log.Printf("\nValues: %v", record.values[0])
}

// This function -- TestValidateSqlite3 -- is a modified https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
// and is therefor licensed under https://github.com/mattn/go-sqlite3/blob/master/LICENSE
//   --- Make sure sqlite3 is working properly ---

func TestValidateSqlite3(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			t.Fatal(err)
		}
	}
	tx.Commit()

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("delete from foo")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		t.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}
