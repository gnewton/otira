package otira

import (
	"errors"
)

const ADDRESS = "address"

func newDefaultTestTable(singleDiscrimField bool) (*TableMeta, error) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		return nil, err
	}

	f0 := new(FieldMetaUint64)
	f0.SetName(pk)
	f0.SetUnique(true)
	err = table.Add(f0)
	if err != nil {
		return nil, err
	}

	f1 := new(FieldMetaString)
	f1.SetName(f_firstname)
	f1.SetFixed(true)
	f1.SetLength(24)
	err = table.Add(f1)
	if err != nil {
		return nil, err
	}

	f2 := new(FieldMetaFloat)
	f2.SetName(f_birth)
	f2.SetLength(32)
	err = table.Add(f2)
	if err != nil {
		return nil, err
	}

	f3 := new(FieldMetaString)
	f3.SetName(f_age)
	f3.SetLength(64)
	err = table.Add(f3)
	if err != nil {
		return nil, err
	}

	f4 := new(FieldMetaFloat)
	f4.SetName(f_height)
	err = table.Add(f4)
	if err != nil {
		return nil, err
	}

	f5 := new(FieldMetaString)
	f5.SetName(tAddress)
	err = table.Add(f5)
	if err != nil {
		return nil, err
	}

	if singleDiscrimField {
		table.SetJoinDiscrimFields(f0)
	} else {
		table.SetJoinDiscrimFields(f0, f1, f2)
	}

	table.SetDone()
	return table, nil
}

var DefaultPID uint64 = 999999

func populateDefaultTableRecord(rec *Record) error {
	err := rec.SetByName(tAddress, "street test")
	if err != nil {
		return err
	}
	err = rec.SetByName(f_firstname, "Robert")
	if err != nil {
		return err
	}

	err = rec.SetByName(f_birth, 10.9)
	if err != nil {
		return err
	}

	err = rec.SetByName(pk, DefaultPID)
	if err != nil {
		return err
	}

	return nil
}

func newOneToManyTestTable() (*TableMeta, *TableMeta, Relation, error) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		return nil, nil, nil, err
	}

	address, err := NewTableMeta(tAddress)
	if err != nil {
		return nil, nil, nil, err
	}

	var pk FieldMeta
	pk = new(FieldMetaUint64)
	pk.SetName("id")
	pk.SetUnique(true)
	err = address.Add(pk)
	if err != nil {
		return nil, nil, nil, err
	}

	var street FieldMeta
	street = new(FieldMetaString)
	street.SetName(fstreet)
	err = address.Add(street)
	if err != nil {
		return nil, nil, nil, err
	}

	var city FieldMeta
	city = new(FieldMetaString)
	city.SetName(fcity)
	err = address.Add(city)
	if err != nil {
		return nil, nil, nil, err
	}

	address.SetDone()

	relation := new(OneToMany)
	relation.name = ADDRESS
	table.AddOneToMany(relation)

	relation.leftTable = table
	relation.rightTable = address

	relation.rightTableUniqueFields = []FieldMeta{city, street}

	relation.leftKeyField = table.GetField(tAddress)
	if relation.leftKeyField == nil {
		return nil, nil, nil, errors.New("Unable to find field [" + tAddress + "] in table [" + table.GetName() + "]")
	}
	relation.rightKeyField = pk

	return table, address, relation, nil
}
