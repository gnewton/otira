package otira

import (
	//"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestCreateSimpleOneToMany(t *testing.T) {
	_, _, _, err := newOneToManyDefaultTables()
	if err != nil {
		t.Fatal(err)
	}

}

func simpleOneToMany() (*Persister, *TableMeta, *TableMeta, *OneToMany, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	//db, err := sql.Open("sqlite3", "smith?_mutex=no&_journal_mode=OFF")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	pers, err := NewPersister(db, NewDialectSqlite3(nil, false))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	address, city, rel, err := newOneToManyDefaultTables()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return pers, address, city, rel, nil
}

func TestVerifySimpleOneToManyCreateWorks(t *testing.T) {
	pers, address, city, _, err := simpleOneToMany()

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
	pers, address, city, _, err := simpleOneToMany()
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

	err = pers.BeginTx()
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

	addressRec1, err := makeAddressRecord1(address)
	err = pers.save(addressRec1)
	if err != nil {
		t.Error(err)
	}

	err = pers.Done()
	if err != nil {
		t.Error(err)
	}

}

func TestVerifySimpleOneToManyInsert_FailMissingCity(t *testing.T) {
	pers, addressTable, city, _, err := simpleOneToMany()
	// defer func() {
	// 	err := pers.Commit()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	err = pers.Done()
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// }()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTable(city)
	if err != nil {
		t.Error(err)
	}
	err = pers.CreateTable(addressTable)
	if err != nil {
		t.Error(err)
	}

	err = pers.BeginTx()
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

	addressRec1, err := makeAddressRecord1(addressTable)
	err = pers.save(addressRec1)
	// Should fail due to foreign key constraints
	if err == nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}

func TestVerifyOneToManyComplexSave(t *testing.T) {
	pers, address, city, one2m, err := simpleOneToMany()
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

	err = pers.BeginTx()
	if err != nil {
		t.Error(err)
	}
	cityRec1, err := makeCityRecord1(city)
	if err != nil {
		t.Error(err)
	}
	addressRec1, err := makeAddressRecord1(address)
	err = addressRec1.AddRelationRecord(one2m, cityRec1)
	if err != nil {
		t.Error(err)
	}
	err = pers.Save(addressRec1)
	if err != nil {
		t.Error(err)
	}
}
