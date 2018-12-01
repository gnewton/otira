package otira

import (
	"testing"
)

func TestJoinKey(t *testing.T) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Error(err)
	}

	table.SetDone()
	rec, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}
	populateDefaultTableRecord(rec)
	t.Log(rec)
	k, err := makeKey(rec)

	t.Log("JOIN KEY key: " + k)
	if err != nil {
		t.Fatal(err)
	}

}

func TestJoinKeyOneDiscrimField(t *testing.T) {
	table, err := newDefaultTestTable(true)
	if err != nil {
		t.Error(err)
	}

	table.SetDone()
	rec, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}
	populateDefaultTableRecord(rec)
	t.Log(rec)
	k, err := makeKey(rec)

	t.Log("JOIN KEY key: " + k)
	if err != nil {
		t.Fatal(err)
	}

}
