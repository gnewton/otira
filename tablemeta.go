package otira

import (
	"errors"
	"log"
)

type TableMeta struct {
	name       string
	fields     []FieldMeta
	fieldsMap  map[string]FieldMeta
	inited     bool
	oneToMany  []*OneToMany
	manyToMany []*ManyToMany
	indexes    []*Index
	done       bool
}

func NewTableMeta(name string) (*TableMeta, error) {
	t := new(TableMeta)
	t.name = name
	t.done = false
	return t, nil
}

func (t *TableMeta) Fields() []FieldMeta {
	return t.fields
}

func (t *TableMeta) GetName() string {
	return t.name
}

func (t *TableMeta) SetDone() {
	t.done = true
}

func (t *TableMeta) Done() bool {
	return t.done
}

func (t *TableMeta) NewRecord() (*Record, error) {
	if !t.done {
		return nil, errors.New("Cannot make new record: TableMeta must be done before using")
	}

	rec := new(Record)
	rec.tableMeta = t
	rec.values = make([]*Field, len(t.fields))
	log.Println("fields", t.fields)
	for i, _ := range rec.values {
		rec.values[i] = new(Field)
		log.Println(i, t.fields[i])
		rec.values[i].fieldMeta = t.fields[i]
	}

	return rec, nil

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

func (t *TableMeta) CreatePreparedStatementInsertAllFields(dialect Dialect) (string, error) {
	a, b := t.CreatePreparedStatementInsertSomeFields(dialect, t.fields...)
	return a, b
}

func (t *TableMeta) CreateTableString(dialect Dialect) (string, error) {
	if !t.done {
		return "", errors.New("Table must be done")
	}
	return dialect.CreateTableString(t), nil
}

func (t *TableMeta) CreatePreparedStatementInsertSomeFields(dialect Dialect, fields ...FieldMeta) (string, error) {
	if !t.done {
		return "", errors.New("Table must be done")
	}

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
	return st, nil
}
