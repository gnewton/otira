package otira

import (
	//"log"
	"testing"
)

func TestCreateTable(t *testing.T) {

	_, err := NewTableMeta("journals")
	if err != nil {
		t.Error(err)
	}

}

const tablename = "journals"
const fieldname0 = "firstname"
const fieldname1 = "lastname"
const fieldname2 = "age"

var fieldnames []string = []string{"city", "stateprovince"}

func TestAddFieldToTable(t *testing.T) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		t.Error(err)
	}
	f := new(FieldMetaString)
	f.SetName(fieldname0)
	t.Log(f.String())
	table.Add(f)
	table.SetDone()
}

func TestAddComplexIndexToTable(t *testing.T) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		t.Error(err)
	}
	f := new(FieldMetaImpl)
	f.SetName("people")
	t.Log(f.String())
	err = table.Add(f)
	if err != nil {
		t.Error(err)
	}
	//TODO
	t.Error(err)

}

func TestCreatePreparedStatementInsert(t *testing.T) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		t.Error(err)
	}
	f := new(FieldMetaImpl)
	f.SetName(fieldname0)
	t.Log(f.String())
	err = table.Add(f)
	if err != nil {
		t.Error(err)
	}
	table.SetDone()
	s := INSERT + tablename + " (" + fieldname0 + ") " + VALUES + " ($1)"
	dialect := new(DialectPostgresql)
	p, err := table.CreatePreparedStatementInsertSomeFields(dialect, f)
	if err != nil {
		t.Error(err)
	}
	if p != s {
		t.Log(s)
		t.Log(p)
		t.Log(p == s)
		t.Error("fdoo")
	}
	//TODO
	//t.Error(err)

}
