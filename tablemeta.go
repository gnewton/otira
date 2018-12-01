package otira

import (
	"errors"
	"log"
	"strconv"
)

type TableMeta struct {
	name          string
	fields        []FieldMeta
	fieldsMap     map[string]FieldMeta
	inited        bool
	oneToMany     []*OneToMany
	oneToManyMap  map[string]*OneToMany
	manyToMany    []*ManyToMany
	manyToManyMap map[string]*ManyToMany
	indexes       []*Index
	done          bool
	ICounter
	discrimFields []FieldMeta
}

func NewTableMeta(name string) (*TableMeta, error) {
	t := new(TableMeta)
	t.name = name
	t.done = false
	return t, nil
}

func (t *TableMeta) SetDiscrimFields(fields ...FieldMeta) {
	t.discrimFields = fields
}

func (t *TableMeta) GetOneToMany(k string) *OneToMany {
	rel, ok := t.oneToManyMap[k]
	if ok {
		return rel
	} else {
		return nil
	}

}

func (t *TableMeta) AddOneToMany(rel *OneToMany) error {
	if rel == nil {
		return errors.New("OneToMany is nil")
	}
	if rel.name == "" {
		return errors.New("Relation cannot have a zero length name")
	}
	log.Println("@@@@@@@@@@@@@@@ " + rel.name)
	if t.oneToMany == nil {
		t.oneToMany = make([]*OneToMany, 0)
		t.oneToManyMap = make(map[string]*OneToMany)
	}
	t.oneToMany = append(t.oneToMany, rel)
	t.oneToManyMap[rel.name] = rel
	return nil
}

func (t *TableMeta) GetField(s string) FieldMeta {
	fm, ok := t.fieldsMap[s]
	if ok {
		return fm
	} else {
		return nil
	}

}

func (t *TableMeta) Fields() []FieldMeta {
	return t.fields
}

func (t *TableMeta) GetName() string {
	return t.name
}

func (t *TableMeta) SetDone() error {

	err := t.validate()
	if err != nil {
		log.Println(err)
		return err
	}
	t.done = true
	return nil
}

// check kto see if there is one and only one primary key
func (t *TableMeta) validate() error {
	if t.fields == nil {
		return errors.New("fields is nil")
	}
	numPrimaryKeys := 0
	for i := 0; i < len(t.fields); i++ {
		if t.fields[i].PrimaryKey() {
			numPrimaryKeys++
		}
	}
	if numPrimaryKeys != 1 {
		return errors.New("Num primary keys != 1; equals:" + strconv.Itoa(numPrimaryKeys))
	}
	return nil
}

func (t *TableMeta) Done() bool {
	return t.done
}

func (t *TableMeta) NewRecordSomeFields(fields ...FieldMeta) (*Record, error) {
	if fields == nil {
		return nil, errors.New("Fields is nil")
	}
	if len(fields) == 0 {
		return nil, errors.New("Fields zero length")
	}
	if !t.done {
		return nil, errors.New("Cannot make new record: TableMeta must be done before using")
	}

	rec, err := newRecord(t, fields, nil)
	if err != nil {
		return nil, err
	}
	//log.Println("fields", fields)
	return rec, nil
}

func (t *TableMeta) NewRecord() (*Record, error) {

	if !t.done {
		return nil, errors.New("Cannot make new record: TableMeta must be done before using")
	}

	return t.NewRecordSomeFields(t.fields...)

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

func (t *TableMeta) createTableString(dialect Dialect) (string, error) {
	if !t.done {
		return "", errors.New("Table must be done")
	}
	if dialect == nil {
		return "", errors.New("Dialect is nil")
	}
	return dialect.CreateTableString(t), nil
}

func (t *TableMeta) CreatePreparedStatementInsertAllFields(dialect Dialect) (string, error) {
	a, b := t.CreatePreparedStatementInsertSomeFields(dialect, t.fields...)
	return a, b
}

func (t *TableMeta) CreatePreparedStatementInsertFromRecord(dialect Dialect, record *Record) (string, error) {
	if record == nil {
		return "", errors.New("Record cannot be nil")
	}
	fields := make([]FieldMeta, 0)
	for i, _ := range record.values {
		fields = append(fields, record.tableMeta.fields[i])
	}
	err := t.SetDone()
	if err != nil {
		return "", err
	}

	return t.CreatePreparedStatementInsertSomeFields(dialect, fields...)
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
