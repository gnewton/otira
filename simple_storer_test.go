package otira

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestSimpleStorer(t *testing.T) {
	//db, err := sql.Open("sqlite3", ":memory:")
	db, err := sql.Open("sqlite3", "foobarMM")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	st := new(SimpleStorer)

	dbopts := DBOptions{
		Dialect:      new(DialectSqlite3),
		TxSize:       10000,
		CreateTables: true,
	}

	var table *TableDef
	table, err = newDefaultTestTable(false)
	if err != nil {
		t.Fatal(err)
	}

	if table == nil {
		t.Fatal(errors.New("Table is nil!!!"))
	}
	err = st.Initialize(db, &dbopts, []*TableDef{table}...)
	if err != nil {
		t.Fatal(err)
	}

	var rec *Record

	for i := 0; i < 10001; i++ {
		rec, err = table.NewRecord()
		if err != nil {
			t.Fatal(err)
		}

		err = populateDefaultTableRecord(rec)
		if err != nil {
			t.Fatal(err)
		}
		err = rec.SetByName(pk, uint64(i))
		if err != nil {
			t.Fatal(err)
		}
		err = st.Save(rec)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = st.Close()
	if err != nil {
		t.Fatal(err)
	}

}
