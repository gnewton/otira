package otira

import (
	"errors"
	"log"
	"testing"
)

func TestJoinKey(t *testing.T) {
	log.Println(1)
	table, err := newDefaultTestTable(false)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(2)
	table.SetDone()
	rec, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(3)
	err = populateDefaultTableRecord(rec)
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatal(err)
	}

	table.SetDone()
	rec, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}
	err = populateDefaultTableRecord(rec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rec)
	k, err := makeKey(rec)

	t.Log("JOIN KEY key: " + k)
	if err != nil {
		t.Fatal(err)
	}

}
func TestJoinCache(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	table, err := newDefaultTestTable(true)
	if err != nil {
		t.Fatal(err)
	}

	table.SetDone()
	rec, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}
	err = populateDefaultTableRecord(rec)
	if err != nil {
		t.Fatal(err)
	}

	jc := NewJoinCache()

	k, exists, err := jc.GetJoinKey(rec)
	if err != nil {
		t.Fatal(err)
	}
	if k != 2 {
		t.Fatal(errors.New("Join key value is incorrect: " + toString(k)))
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
		t.Fatal(err)
	}
	table.useRecordPrimaryKeys = true
	table.SetDone()
	rec, err := table.NewRecord()
	if err != nil {
		t.Fatal(err)
	}
	err = populateDefaultTableRecord(rec)
	if err != nil {
		t.Fatal(err)
	}

	jc := NewJoinCache()
	log.Println(rec)
	k, exists, err := jc.GetJoinKey(rec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(k)
	t.Log(exists)
	if k != DefaultPID {
		t.Fatal(errors.New("Join key value is incorrect: " + toString(k)))
	}
	if exists {
		t.Fatal(errors.New("Join key should NOT be in the cache"))
	}

	k, exists, err = jc.GetJoinKey(rec)
	if err != nil {
		t.Fatal(err)
	}
	if k != DefaultPID {
		t.Fatal(errors.New("Join key value is incorrect" + toString(k)))
	}
	if !exists {
		t.Fatal(errors.New("Join key SHOULD be in the cache"))
	}
	t.Log(k)
	t.Log(exists)
}
