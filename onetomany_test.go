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

	err = pers.CreateTables(city, address)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifysShallowOneToManyInsert(t *testing.T) {
	pers, address, city, _, err := simpleOneToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(city, address)
	if err != nil {
		t.Fatal(err)
	}

	cityRec1, err := makeCityRecord1(city, City1PK)
	if err != nil {
		t.Fatal(err)
	}
	// Here we are individually saving, calling save() not Save()
	err = pers.save(cityRec1)
	if err != nil {
		t.Fatal(err)
	}

	addressRec1, err := makeAddressRecord1(address, Address1PK)
	err = pers.save(addressRec1)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.Done()
	if err != nil {
		t.Fatal(err)
	}

}

func TestVerifysDeepOneToManyInsert(t *testing.T) {
	//func Foo(t *testing.T) {
	pers, address, city, _, err := simpleOneToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(city, address)
	if err != nil {
		t.Fatal(err)
	}

	cityRec1, err := makeCityRecord1(city, City1PK)
	if err != nil {
		t.Fatal(err)
	}

	// Here we are deep saving, calling Save() not save();
	//   Save() should recursively save City, then Address, populating address.cityfk with city.fk
	addressRec1, err := makeAddressRecord1(address, Address1PK)
	err = addressRec1.AddRelationRecord(nil, cityRec1)
	if err == nil {
		t.Fatal(err)
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

	err = pers.CreateTables(city, addressTable)
	if err != nil {
		t.Fatal(err)
	}
	cityRec2, err := makeCityRecord2(city)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.save(cityRec2)
	if err != nil {
		t.Fatal(err)
	}

	addressRec1, err := makeAddressRecord1(addressTable, Address1PK)
	err = pers.save(addressRec1)
	// Should fail due to foreign key constraints
	if err == nil {
		t.Fatal(err)
	} else {
		t.Log(err)
	}
}

func TestVerifyOneToManyComplexSaveJoinCache(t *testing.T) {
	pers, address, city, one2m, err := simpleOneToMany()
	city.useRecordPrimaryKeys = true
	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(city, address)
	if err != nil {
		t.Fatal(err)
	}
	cityRec1, err := makeCityRecord1(city, City1PK)
	if err != nil {
		t.Fatal(err)
	}
	addressRec1, err := makeAddressRecord1(address, Address1PK)
	err = addressRec1.AddRelationRecord(one2m, cityRec1)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.Save(addressRec1)
	if err != nil {
		t.Fatal(err)
	}
	cityRec2, err := makeCityRecord1(city, City1PK)
	if err != nil {
		t.Fatal(err)
	}
	addressRec2, err := makeAddressRecord1(address, Address2PK)
	err = addressRec2.AddRelationRecord(one2m, cityRec2)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.Save(addressRec2)
	////

	defer pers.Done()
}

func TestVerifyOneToManyComplexSave(t *testing.T) {
	pers, address, city, one2m, err := simpleOneToMany()
	city.useRecordPrimaryKeys = true
	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(city, address)
	if err != nil {
		t.Fatal(err)
	}
	cityRec1, err := makeCityRecord1(city, City1PK)
	if err != nil {
		t.Fatal(err)
	}
	addressRec1, err := makeAddressRecord1(address, Address1PK)
	err = addressRec1.AddRelationRecord(one2m, cityRec1)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.Save(addressRec1)
	if err != nil {
		t.Fatal(err)
	}
	defer pers.Done()
}
