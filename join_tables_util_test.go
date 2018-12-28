package otira

import (
	"errors"
)

const STREET = "street"
const CITY = "city"
const CITYFK = "cityfk"
const NAME = "name"

func newOneToManyDefaultTables() (*TableMeta, *TableMeta, *OneToMany, error) {
	addressTable, err := makeAddressTable()
	if err != nil {
		return nil, nil, nil, err
	}
	cityTable, err := makeCityTable()
	if err != nil {
		return nil, nil, nil, err
	}
	one2m := NewOneToMany()
	one2m.leftTable = addressTable
	one2m.rightTable = cityTable
	cityField := addressTable.GetField(CITYFK)
	if cityField == nil {
		return nil, nil, nil, errors.New("Cannot find field: " + CITYFK + " in table: " + addressTable.name)

	}
	one2m.leftKeyField = cityField
	one2m.rightKeyField = cityTable.PrimaryKey()
	one2m.rightTableUniqueFields = []FieldMeta{cityField}

	addressTable.AddOneToMany(one2m)
	addressTable.SetDone()
	return addressTable, cityTable, one2m, nil
}

func makeAddressTable() (*TableMeta, error) {
	addressTable, err := NewTableMeta(ADDRESS)
	if err != nil {
		return nil, err
	}

	id := new(FieldMetaUint64)
	id.SetName(pk)
	id.SetUnique(true)
	err = addressTable.Add(id)
	if err != nil {
		return nil, err
	}

	streetField := new(FieldMetaString)
	streetField.SetName(STREET)
	streetField.SetFixed(true)
	streetField.SetLength(24)
	err = addressTable.Add(streetField)
	if err != nil {
		return nil, err
	}

	cityField := new(FieldMetaInt)
	cityField.SetName(CITYFK)
	err = addressTable.Add(cityField)
	if err != nil {
		return nil, err
	}

	return addressTable, nil
}

func makeCityTable() (*TableMeta, error) {
	cityTable, err := NewTableMeta(CITY)
	if err != nil {
		return nil, err
	}
	id := new(FieldMetaUint64)
	id.SetName(pk)
	id.SetUnique(true)
	err = cityTable.Add(id)
	if err != nil {
		return nil, err
	}

	nameField := new(FieldMetaString)
	nameField.SetName(NAME)
	nameField.SetFixed(true)
	nameField.SetLength(24)
	err = cityTable.Add(nameField)
	if err != nil {
		return nil, err
	}
	cityTable.SetDiscrimFields(nameField)
	cityTable.SetDone()
	return cityTable, err
}

const City1PK = uint64(42)
const City1Name = "New York"
const City2PK = 73
const City2Name = "Montreal"

const Address1PK = uint64(675)
const Address1Street = "Main St."
const Address2PK = uint64(88908)

func makeCityRecord1(t *TableMeta, citypkvalue uint64) (*Record, error) {
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

func makeCityRecord2(t *TableMeta) (*Record, error) {
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

func makeAddressRecord1(t *TableMeta, addressPkValue uint64) (*Record, error) {
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
