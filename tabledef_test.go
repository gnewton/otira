package otira

import (
	//"log"
	"testing"
)

func TestCreateTable(t *testing.T) {

	_, err := NewTableDef("journals")
	if err != nil {
		t.Fatal(err)
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
	table, err := NewTableDef(tablename)
	if err != nil {
		t.Fatal(err)
	}
	f := new(FieldDefUint64)
	f.SetName(f_firstname)
	t.Log(f.String())
	err = table.Add(f)
	if err != nil {
		t.Fatal(err)
	}
	table.SetDone()
}

func TestAddComplexIndexToTable(t *testing.T) {
	table, err := NewTableDef(tablename)
	if err != nil {
		t.Fatal(err)
	}
	f := new(FieldDefUint64)
	f.SetName("people")

	t.Log(f.String())
	err = table.Add(f)
	if err != nil {
		t.Fatal(err)
	}
	//TODO
	//t.Fatal(err)

}

func TestCreatePreparedStatementInsert(t *testing.T) {
	table, err := NewTableDef(tablename)
	if err != nil {
		t.Fatal(err)
	}

	pk := new(FieldDefUint64)
	pk.SetName("id")

	t.Log(pk.String())
	err = table.Add(pk)
	f := new(FieldDefString)
	f.SetName(f_firstname)
	t.Log(f.String())
	err = table.Add(f)
	if err != nil {
		t.Fatal(err)
	}
	err = table.SetDone()
	if err != nil {
		t.Fatal(err)
	}
	s := INSERT + tablename + " (id, " + f_firstname + ") " + VALUES + " ($1, $2)"
	dialect := new(DialectPostgresql)
	p, err := table.CreatePreparedStatementInsertSomeFields(dialect, pk, f)
	if err != nil {
		t.Fatal(err)
	}
	if p != s {
		t.Log(s)
		t.Log(p)
		t.Log(p == s)
		t.Fatal("fdoo")
	}
	//TODO
	//t.Fatal(err)

}

func TestAddTwoFieldsWithSameName(t *testing.T) {
	table, err := NewTableDef(tablename)
	if err != nil {
		t.Fatal(err)
	}
	f := new(FieldDefUint64)
	f.SetName(f_firstname)
	t.Log(f.String())
	err = table.Add(f)
	if err != nil {
		t.Fatal(err)
	}

	f2 := new(FieldDefUint64)
	f2.SetName("fred")
	err = table.Add(f2)
	if err != nil {
		t.Fatal(err)
	}

	f3 := new(FieldDefUint64)
	f3.SetName("fred")
	err = table.Add(f3)
	if err == nil {
		t.Fatal(err)
	}
}
