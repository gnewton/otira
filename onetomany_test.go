package otira

import (
	//"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestCreateSimpleOneToMany(t *testing.T) {
	_, _, err := newOneToManyDefaultTables()
	if err != nil {
		t.Fatal(err)
	}

}

func simpleOneToMany() (*Persister, *TableMeta, *TableMeta, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, nil, nil, err
	}
	pers, err := NewPersister(db, new(DialectSqlite3))
	if err != nil {
		return nil, nil, nil, err
	}

	address, city, err := newOneToManyDefaultTables()
	if err != nil {
		return nil, nil, nil, err
	}

	return pers, address, city, nil
}

func TestVerifySimpleOneToManyCreateWorks(t *testing.T) {
	pers, address, city, err := simpleOneToMany()
	defer pers.Done()
	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTable(city)
	if err != nil {
		t.Error(err)
	}

	err = pers.CreateTable(address)
	if err != nil {
		t.Error(err)
	}

}

func TestVerifySimpleOneToManyInsert(t *testing.T) {
	pers, address, city, err := simpleOneToMany()
	defer pers.Done()
	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTable(city)
	if err != nil {
		t.Error(err)
	}

	cityRec1, err := makeCityRecord1(city)
	if err != nil {
		t.Error(err)
	}
	err = pers.save(cityRec1)
	if err != nil {
		t.Error(err)
	}

	err = pers.CreateTable(address)
	if err != nil {
		t.Error(err)
	}
	addressRec1, err := makeAddressRecord1(address)
	err = pers.save(addressRec1)
	if err != nil {
		t.Error(err)
	}
}

func TestVerifySimpleOneToManyInsert_FailMissingCity(t *testing.T) {
	//pers, address, city, err := simpleOneToMany()
	pers, address, city, err := simpleOneToMany()
	defer pers.Done()
	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTable(city)
	if err != nil {
		t.Error(err)
	}
	cityRec2, err := makeCityRecord2(city)
	if err != nil {
		t.Error(err)
	}
	err = pers.save(cityRec2)
	if err != nil {
		t.Error(err)
	}

	err = pers.CreateTable(address)
	if err != nil {
		t.Error(err)
	}
	addressRec1, err := makeAddressRecord1(address)
	err = pers.save(addressRec1)
	if err != nil {
		t.Error(err)
	}
}
