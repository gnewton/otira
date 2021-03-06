package otira

import (
	"errors"
)

const STREET = "street"
const CITY = "city"
const CITYFK = "cityfk"
const NAME = "name"

func newOneToManyDefaultTables() (*TableDef, *TableDef, *OneToMany, error) {
	addressTable, err := makeAddressTable()
	if err != nil {
		return nil, nil, nil, err
	}
	cityTable, err := makeCityTable()
	if err != nil {
		return nil, nil, nil, err
	}
	one2m := NewOneToMany()
	one2m.LeftTable = addressTable
	one2m.RightTable = cityTable
	cityField := addressTable.GetField(CITYFK)
	if cityField == nil {
		return nil, nil, nil, errors.New("Cannot find field: " + CITYFK + " in table: " + addressTable.name)

	}
	one2m.LeftKeyField = cityField
	one2m.RightKeyField = cityTable.PrimaryKey()
	one2m.RightTableUniqueFields = []FieldDef{cityField}

	addressTable.AddOneToMany(one2m)
	addressTable.SetDone()
	return addressTable, cityTable, one2m, nil
}

func makeAddressTable() (*TableDef, error) {
	addressTable, err := NewTableDef(ADDRESS)
	if err != nil {
		return nil, err
	}

	id := new(FieldDefInt64)
	id.SetName(pk)
	id.SetUnique(true)
	err = addressTable.Add(id)
	if err != nil {
		return nil, err
	}

	streetField := new(FieldDefString)
	streetField.SetName(STREET)
	streetField.SetFixed(true)
	streetField.SetLength(24)
	err = addressTable.Add(streetField)
	if err != nil {
		return nil, err
	}

	cityField := new(FieldDefInt)
	cityField.SetName(CITYFK)
	err = addressTable.Add(cityField)
	if err != nil {
		return nil, err
	}

	return addressTable, nil
}

func makeCityTable() (*TableDef, error) {
	cityTable, err := NewTableDef(CITY)
	if err != nil {
		return nil, err
	}
	id := new(FieldDefInt64)
	id.SetName(pk)
	id.SetUnique(true)
	err = cityTable.Add(id)
	if err != nil {
		return nil, err
	}

	nameField := new(FieldDefString)
	nameField.SetName(NAME)
	nameField.SetFixed(true)
	nameField.SetLength(24)
	err = cityTable.Add(nameField)
	if err != nil {
		return nil, err
	}
	cityTable.SetJoinDiscrimFields(nameField)
	cityTable.SetDone()
	return cityTable, err
}

const City1PK = int64(42)
const City1Name = "New York"
const City2PK = int64(73)
const City2Name = "Montreal"

const Address1PK = int64(675)
const Address1Street = "Main St."
const Address2PK = int64(88908)

func makeCityRecord1(t *TableDef, citypkvalue int64) (*Record, error) {
	rec, err := t.NewRecord()
	if err != nil {
		return nil, err
	}
	err = rec.SetByName(pk, citypkvalue)
	if err != nil {
		return nil, err
	}

	err = rec.SetByName(NAME, City1Name)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func makeCityRecord2(t *TableDef) (*Record, error) {
	rec, err := t.NewRecord()
	if err != nil {
		return nil, err
	}
	err = rec.SetByName(pk, City2PK)
	if err != nil {
		return nil, err
	}

	err = rec.SetByName(NAME, City2Name)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func makeAddressRecord1(t *TableDef, addressPkValue int64) (*Record, error) {
	rec, err := t.NewRecord()
	if err != nil {
		return nil, err
	}
	err = rec.SetByName(pk, addressPkValue)
	if err != nil {
		return nil, err
	}

	err = rec.SetByName(STREET, Address1Street)
	if err != nil {
		return nil, err
	}
	err = rec.SetByName(CITYFK, City1PK)
	if err != nil {
		return nil, err
	}

	return rec, nil
}
