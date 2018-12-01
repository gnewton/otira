package otira

import (
	"errors"
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
func TestJoinCache(t *testing.T) {
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

	jc := NewJoinCache()

	k, exists, err := jc.GetJoinKey(rec)
	if err != nil {
		t.Fatal(err)
	}
	if k != 1 {
		t.Fatal(errors.New("Join key value is incorrect"))
	}
	if exists {
		t.Fatal(errors.New("Join key should NOT be in the cache"))
	}
	t.Log(k)
	t.Log(exists)
}

func TestJoinCacheDuplicate(t *testing.T) {
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

	jc := NewJoinCache()

	k, exists, err := jc.GetJoinKey(rec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(k)
	t.Log(exists)
	if k != 1 {
		t.Fatal(errors.New("Join key value is incorrect" + toString(k)))
	}
	if exists {
		t.Fatal(errors.New("Join key should NOT be in the cache"))
	}

	k, exists, err = jc.GetJoinKey(rec)
	if err != nil {
		t.Fatal(err)
	}
	if k != 1 {
		t.Fatal(errors.New("Join key value is incorrect" + toString(k)))
	}
	if !exists {
		t.Fatal(errors.New("Join key SHOULD be in the cache"))
	}
	t.Log(k)
	t.Log(exists)
}
