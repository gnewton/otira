package otira

func defaultTestTable() (*TableMeta, error) {
	table, err := NewTableMeta(tablename)
	if err != nil {
		return nil, err
	}

	f1 := new(FieldMetaString)
	f1.SetName(fieldname0)
	f1.SetFixed(true)
	table.Add(f1)

	f2 := new(FieldMetaInt)
	f2.SetName(fieldname2)
	table.Add(f2)

	table.SetDone()
	return table, nil
}
