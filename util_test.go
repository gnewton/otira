package otira

import (
	"errors"
)

const ADDRESS = "ADDRESS"

func newDefaultTestTable() (*TableMeta, error) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		return nil, err
	}

	f0 := new(FieldMetaInt)
	f0.SetName(pk)
	f0.SetPrimaryKey(true)
	f0.SetUnique(true)
	table.Add(f0)

	f1 := new(FieldMetaString)
	f1.SetName(fieldname0)
	f1.SetFixed(true)
	f1.SetLength(24)
	table.Add(f1)

	f2 := new(FieldMetaFloat)
	f2.SetName(fieldname2)
	f2.SetLength(32)
	table.Add(f2)

	f3 := new(FieldMetaString)
	f3.SetName(fieldname1)
	f3.SetLength(64)
	table.Add(f3)

	f4 := new(FieldMetaFloat)
	f4.SetName(fieldname3)
	table.Add(f4)

	f5 := new(FieldMetaInt)
	f5.SetName(tAddress)
	table.Add(f5)

	table.SetDone()
	return table, nil
}

func newOneToManyTestTable() (*TableMeta, *TableMeta, Relation, error) {
	table, err := newDefaultTestTable()
	if err != nil {
		return nil, nil, nil, err
	}

	address, err := NewTableMeta(tAddress)
	if err != nil {
		return nil, nil, nil, err
	}

	var street FieldMeta
	street = new(FieldMetaString)
	street.SetName(fstreet)
	address.Add(street)

	var city FieldMeta
	city = new(FieldMetaString)
	city.SetName(fcity)
	address.Add(city)

	var pk FieldMeta
	pk = new(FieldMetaInt)
	pk.SetName("id")
	pk.SetPrimaryKey(true)
	pk.SetUnique(true)
	address.Add(pk)
	address.SetDone()

	relation := new(OneToMany)
	relation.name = ADDRESS
	err = table.AddOneToMany(relation)
	if err != nil {
		return nil, nil, nil, err
	}
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
