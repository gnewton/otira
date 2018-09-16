package otira

func defaultTestTable() (*TableMeta, error) {
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

	f2 := new(FieldMetaInt)
	f2.SetName(fieldname2)
	f2.SetLength(32)
	table.Add(f2)

	f3 := new(FieldMetaString)
	f3.SetName(fieldname1)
	f3.SetLength(64)
	table.Add(f3)

	table.SetDone()
	return table, nil
}
