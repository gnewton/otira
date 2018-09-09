package otira

import (
	"testing"
)

func TestCreateTable(t *testing.T) {

	_, err := NewTableMeta("journals")
	if err != nil {
		t.Error(err)
	}

}

func TestAddFieldToTable(t *testing.T) {

	table, err := NewTableMeta("journals")
	if err != nil {
		t.Error(err)
	}
	table.Add(new(FieldMetaImpl))

}
