package otira

import (
	"errors"
)

const ADDRESS = "address"

func newDefaultTestTable(singleDiscrimField bool) (*TableDef, error) {
	table, err := NewTableDef(tablename)
	if err != nil {
		return nil, err
	}

	f0 := new(FieldDefInt64)
	f0.SetName(pk)
	f0.SetUnique(true)
	err = table.Add(f0)
	if err != nil {
		return nil, err
	}

	f1 := new(FieldDefString)
	f1.SetName(f_firstname)
	f1.SetFixed(true)
	f1.SetLength(24)
	err = table.Add(f1)
	if err != nil {
		return nil, err
	}

	f2 := new(FieldDefFloat)
	f2.SetName(f_birth)
	f2.SetLength(32)
	err = table.Add(f2)
	if err != nil {
		return nil, err
	}

	f3 := new(FieldDefString)
	f3.SetName(f_age)
	f3.SetLength(64)
	err = table.Add(f3)
	if err != nil {
		return nil, err
	}

	f4 := new(FieldDefFloat)
	f4.SetName(f_height)
	err = table.Add(f4)
	if err != nil {
		return nil, err
	}

	f5 := new(FieldDefString)
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

var DefaultPID int64 = 999999

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

func newOneToManyTestTable() (*TableDef, *TableDef, Relation, error) {
	table, err := newDefaultTestTable(false)
	if err != nil {
		return nil, nil, nil, err
	}

	address, err := NewTableDef(tAddress)
	if err != nil {
		return nil, nil, nil, err
	}

	var pk FieldDef
	pk = new(FieldDefInt64)
	pk.SetName("id")
	pk.SetUnique(true)
	err = address.Add(pk)
	if err != nil {
		return nil, nil, nil, err
	}

	var street FieldDef
	street = new(FieldDefString)
	street.SetName(fstreet)
	err = address.Add(street)
	if err != nil {
		return nil, nil, nil, err
	}

	var city FieldDef
	city = new(FieldDefString)
	city.SetName(fcity)
	err = address.Add(city)
	if err != nil {
		return nil, nil, nil, err
	}

	address.SetDone()

	relation := new(OneToMany)
	relation.name = ADDRESS
	table.AddOneToMany(relation)

	relation.LeftTable = table
	relation.RightTable = address

	relation.RightTableUniqueFields = []FieldDef{city, street}

	relation.LeftKeyField = table.GetField(tAddress)
	if relation.LeftKeyField == nil {
		return nil, nil, nil, errors.New("Unable to find field [" + tAddress + "] in table [" + table.Name() + "]")
	}
	relation.RightKeyField = pk

	return table, address, relation, nil
}
