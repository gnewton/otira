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

var fieldnames []string = []string{"city", "stateprovince"}

func TestAddFieldToTable(t *testing.T) {

	table, err := NewTableMeta("journals")
	if err != nil {
		t.Error(err)
	}
	f := new(FieldMetaImpl)
	f.SetName(fieldname0)
	t.Log(f.String())
	table.Add(f)
}

func TestAddComplexIndexToTable(t *testing.T) {
	table, err := NewTableMeta("journals")
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
	table, err := NewTableMeta("journals")
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
	t.Log(table.CreatePreparedStatementInsertSomeFields("postgresql", f))
	//TODO
	t.Error(err)

}
