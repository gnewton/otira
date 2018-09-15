package otira

type TableMeta struct {
	name       string
	fields     []FieldMeta
	fieldsMap  map[string]FieldMeta
	inited     bool
	oneToMany  []*OneToMany
	manyToMany []*ManyToMany
	indexes    []*Index
}

func NewTableMeta(name string) (*TableMeta, error) {
	t := new(TableMeta)
	t.name = name
	return t, nil
}

func (t *TableMeta) AddIndex(name string, field0, field1 *FieldMeta, fields ...*FieldMeta) {
	if t.indexes == nil {
		t.indexes = make([]*Index, 1)
	}
	index := NewIndex(name, field0, field1, fields...)
	t.indexes = append(t.indexes, index)

}

func (t *TableMeta) Add(f FieldMeta) error {
	if err := baseFieldMetaErrors(f); err != nil {
		return err
	}

	if t.fieldsMap == nil {
		t.fieldsMap = make(map[string]FieldMeta, 0)
	}

	if t.fields == nil {
		t.fields = make([]FieldMeta, 0)
	}

	t.fieldsMap[f.Name()] = f
	t.fields = append(t.fields, f)

	return nil
}

func (t *TableMeta) CreatePreparedStatementInsert(dialect string) string {
	return t.CreatePreparedStatementInsertSomeFields(dialect, t.fields...)
}

func (t *TableMeta) CreatePreparedStatementInsertSomeFields(dialect string, fields ...FieldMeta) string {

	st := "INSERT INTO " + t.name + " ("
	values := "("
	for i, _ := range fields {
		if i != 0 {
			st += ", "
			values += ", "
		}
		st += fields[i].Name()
		values += preparedValueFormat(dialect, i)
	}

	st += ")"
	values += ")"

	st = st + " VALUES " + values
	return st

	// //+"(column names) VALUES  (?)"

	// switch dialect {
	// case "oracle":

	// case "mysql":
	// case "sqlite3":

	// case "postgresql":
	// }

	// return "TODO FIXXX"
}
