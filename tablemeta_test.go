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

const tablename = "person"
const f_firstname = "firstname"
const f_birth = "birth"
const f_age = "age"
const f_height = "height"
const tAddress = "address"
const fstreet = "street"
const fcity = "city"
const pk = "id"

var fieldnames []string = []string{"city", "stateprovince"}

func TestAddFieldToTable(t *testing.T) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		t.Error(err)
	}
	f := new(FieldMetaString)
	f.SetName(f_firstname)
	t.Log(f.String())
	table.Add(f)
	table.SetPrimaryKey(f)
	table.SetDone()
}

func TestAddComplexIndexToTable(t *testing.T) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		t.Error(err)
	}
	f := new(FieldMetaString)
	f.SetName("people")
	table.SetPrimaryKey(f)
	t.Log(f.String())
	err = table.Add(f)
	if err != nil {
		t.Error(err)
	}
	//TODO
	//t.Error(err)

}

func TestCreatePreparedStatementInsert(t *testing.T) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		t.Error(err)
	}

	pk := new(FieldMetaInt)
	pk.SetName("id")
	pk.SetPrimaryKey(true)
	table.SetPrimaryKey(pk)
	t.Log(pk.String())
	err = table.Add(pk)
	f := new(FieldMetaString)
	f.SetName(f_firstname)
	t.Log(f.String())
	err = table.Add(f)
	if err != nil {
		t.Error(err)
	}
	err = table.SetDone()
	if err != nil {
		t.Error(err)
	}
	s := INSERT + tablename + " (id, " + f_firstname + ") " + VALUES + " ($1, $2)"
	dialect := new(DialectPostgresql)
	p, err := table.CreatePreparedStatementInsertSomeFields(dialect, pk, f)
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

func TestDetectNoPrimaryKeyWithValidate(t *testing.T) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		t.Error(err)
	}

	f0 := new(FieldMetaInt)
	f0.SetName(pk)
	f0.SetUnique(true)
	table.Add(f0)

	f1 := new(FieldMetaString)
	f1.SetName(f_firstname)
	f1.SetFixed(true)
	f1.SetLength(24)
	table.Add(f1)

	err = table.SetDone()
	if err == nil {
		t.Error(err)
	}
}
