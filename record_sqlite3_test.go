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

// This function -- TestValidateSqlite3 -- is a slightly modified https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
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
