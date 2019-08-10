package otira

import (
	"io/ioutil"
	"testing"
)

func TestOpenHC(t *testing.T) {
	hc, err := NewBadgerSet("/tmp/tmpbadger")

	if err != nil {
		t.Fatal(err)
	}

	err = hc.Close()
	if err != nil {
		t.Fatal(err)
	}

}

func TestAddHC(t *testing.T) {
	dir, err := ioutil.TempDir("", "otira_badger_")
	if err != nil {
		t.Fatal(err)
	}
	hc, err := NewBadgerSet(dir)

	if err != nil {
		t.Fatal(err)
	}

	var v int64 = 100
	hc.Put(v)

	v = int64(199)
	hc.Put(v)

	if ok, err := hc.Contains(v); !ok {
		t.Fatal(err)
	}

	err = hc.Close()
	if err != nil {
		t.Fatal(err)
	}

}
