package otira

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// func TestMain(m *testing.M) {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)
// 	// call flag.Parse() here if TestMain uses flags
// 	os.Exit(m.Run())
// }

func TestNewPreparedStatement(t *testing.T) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Fatal(err)
	}

	prep, err := table.CreatePreparedStatementInsertAllFields(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(prep)
}

func TestNewTableOneToMany(t *testing.T) {
	_, _, _, err := newOneToManyTestTable()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddOneToManyRecordToMainRecord(t *testing.T) {

	table, person, relation, err := newOneToManyTestTable()
	if err != nil {
		t.Fatal(err)
	}

	tableRecord, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}

	personRecord, err := person.NewRecord()
	if err != nil {
		t.Fatal(err)
	}

	err = tableRecord.AddRelationRecord(relation, personRecord)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddOneToManyRecordToMainRecordByRelationStringName_BadName(t *testing.T) {

	table, _, _, err := newOneToManyTestTable()
	if err != nil {
		t.Fatal(err)
	}

	relation := table.GetOneToMany(ADDRESS + "m")
	if relation != nil {
		t.Fatal(errors.New("Cannot find the relation by string key"))
	}
}

func TestTableCreate(t *testing.T) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Fatal(err)
	}

	cr, err := table.createTableString(new(DialectSqlite3))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cr)

}

func TestCreateTableSyntax(t *testing.T) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	create, err := table.createTableString(new(DialectSqlite3))
	_, err = db.Exec(create)

	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateTableSyntaxFail(t *testing.T) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	dialect := new(DialectSqlite3)
	create, err := dialect.CreateTableString(table)
	_, err = db.Exec("ZZ " + create)

	if err == nil {
		t.Fatal(err)
	}

}

func TestWriteRecordsFromTableMeta(t *testing.T) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Fatal(err)
	}
	//record, err := table.NewRecordSomeFields(table.fields[0], table.fields[1], table.fields[2])
	record, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}
	record.validating = true

	if err != nil {
		t.Fatal(err)
	}
	err = record.Set(0, 42)
	v, err := record.tableMeta.Next()
	if err != nil {
		t.Fatal(err)
	}
	err = record.Set(0, v)
	if err != nil {
		t.Fatal(err)
	}

	err = record.Set(1, "mm")
	if err != nil {
		t.Fatal(err)
	}

	err = record.Set(2, 4.5)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dialect := new(DialectSqlite3)
	pers, err := NewPersister(db, dialect)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(table)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.save(record)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	for i := 100; i < 201; i++ {
		record.Reset()
		record.Set(0, i)
		err = pers.save(record)
		if err != nil {
			t.Log(i)
			t.Fatal(err)
		}
	}

	//t.Log(result.RowsAffected())
	pers.commit()

}

func TestRecordFromTableMetaTODO(t *testing.T) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Fatal(err)
	}
	record, err := table.NewRecord()

	err = record.Set(0, "mm")
	if err != nil {
		t.Fatal(err)
	}
	v := 44
	err = record.Set(1, v)
	if err != nil {
		t.Fatal(err)
	}
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
