package otira

import (
	"errors"
	"log"
)

//Assumes 0th field is primary key, type FieldMetaUint64
type TableDef struct {
	ICounter
	UseRecordPrimaryKeys bool
	created              bool
	done                 bool
	fields               []FieldMeta
	fieldsMap            map[string]FieldMeta
	indexes              []*Index
	isJoinTable          bool
	inited               bool
	joinDiscrimFields    []FieldMeta
	manyToMany           []*ManyToMany
	manyToManyMap        map[string]*ManyToMany
	name                 string
	oneToMany            []*OneToMany
	oneToManyMap         map[string]*OneToMany
	primaryKey           FieldMeta
	primaryKeyIndex      int
	validating           bool
	writeCache           Set
}

func NewTableDef(name string) (*TableDef, error) {
	t := new(TableDef)
	t.name = name
	t.done = false
	t.created = false
	t.UseRecordPrimaryKeys = false
	t.counter = 0
	t.validating = true
	t.isJoinTable = false

	t.oneToMany = make([]*OneToMany, 0)
	t.oneToManyMap = make(map[string]*OneToMany)
	t.manyToMany = make([]*ManyToMany, 0)
	t.manyToManyMap = make(map[string]*ManyToMany)
	t.indexes = make([]*Index, 0)
	t.fieldsMap = make(map[string]FieldMeta, 0)
	t.fields = make([]FieldMeta, 0)
	var err error
	t.writeCache, err = NewBadgerSet("./foom")
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TableDef) AddOneToMany(one2m *OneToMany) {
	t.oneToMany = append(t.oneToMany, one2m)
	t.oneToManyMap[one2m.name] = one2m
}

func (t *TableDef) AddManyToMany(m2m *ManyToMany) {
	t.manyToMany = append(t.manyToMany, m2m)
	t.manyToManyMap[m2m.name] = m2m
	makeM2MJoinTable(m2m)
}

func makeM2MJoinTable(m2m *ManyToMany) error {
	joinTable, err := NewTableDef(m2m.leftTable.name + "_" + m2m.rightTable.name)
	if err != nil {
		return err
	}
	joinTable.isJoinTable = true
	m2m.joinTable = joinTable

	left := new(FieldMetaUint64)
	left.SetName(m2m.leftTable.name)
	err = joinTable.Add(left)
	if err != nil {
		return err
	}

	right := new(FieldMetaUint64)
	right.SetName(m2m.rightTable.name)
	err = joinTable.Add(right)
	if err != nil {
		return err
	}
	err = joinTable.SetDone()
	return err
}

func (t *TableDef) PrimaryKey() FieldMeta {
	return t.fields[0]
}

func (t *TableDef) SetJoinDiscrimFields(fields ...FieldMeta) {
	t.joinDiscrimFields = fields
}

func (t *TableDef) GetOneToMany(k string) *OneToMany {
	rel, ok := t.oneToManyMap[k]
	if ok {
		return rel
	} else {
		return nil
	}

}

func (t *TableDef) AddOneToMany_OLD(rel *OneToMany) error {
	if rel == nil {
		return errors.New("OneToMany is nil")
	}
	if rel.name == "" {
		return errors.New("Relation cannot have a zero length name")
	}
	log.Println("@@@@@@@@@@@@@@@ " + rel.name)
	t.oneToMany = append(t.oneToMany, rel)
	t.oneToManyMap[rel.name] = rel
	return nil
}

func (t *TableDef) GetField(s string) FieldMeta {
	fm, ok := t.fieldsMap[s]
	if ok {
		return fm
	} else {
		return nil
	}

}

func (t *TableDef) Fields() []FieldMeta {
	return t.fields
}

func (t *TableDef) GetName() string {
	return t.name
}

func (t *TableDef) SetDone() error {
	err := t.validate()
	if err != nil {
		log.Println(err)
		return err
	}
	t.done = true
	return nil
}

// check kto see if there is one and only one primary key
func (t *TableDef) validate() error {
	if t == nil {
		return errors.New("FATAL !!!!!!!!!!!!!!!!!!! tablemeta is nil")
	}
	if t.fields == nil {
		return errors.New("fields is nil")
	}
	return nil
}

func (t *TableDef) Done() bool {
	return t.done
}

func (t *TableDef) NewRecordSomeFields(fields ...FieldMeta) (*Record, error) {
	if fields == nil {
		return nil, errors.New("Fields is nil")
	}
	if len(fields) == 0 {
		return nil, errors.New("Fields zero length")
	}
	if !t.done {
		return nil, errors.New("Cannot make new record: TableDef must be done before using")
	}

	rec, err := newRecord(t, fields, nil)
	if err != nil {
		return nil, err
	}
	//log.Println("fields", fields)
	if !t.UseRecordPrimaryKeys {
		pk, err := t.Next()
		if err != nil {
			return nil, err
		}
		rec.SetByName(t.PrimaryKey().Name(), pk)
	}
	rec.SetByName(t.PrimaryKey().Name(), -1)
	return rec, nil
}

func (t *TableDef) NewRecord() (*Record, error) {

	if !t.done {
		return nil, errors.New("Cannot make new record: TableDef must be done before using")
	}

	return t.NewRecordSomeFields(t.fields...)

}

func (t *TableDef) AddIndex(name string, field0, field1 *FieldMeta, fields ...*FieldMeta) error {
	if t.indexes == nil {
		return errors.New("Index is nil")
	}
	index := NewIndex(name, field0, field1, fields...)
	t.indexes = append(t.indexes, index)
	return nil

}

func (t *TableDef) Add(f FieldMeta) error {
	if err := baseFieldMetaErrors(f); err != nil {
		return err
	}

	if t.fieldsMap == nil {
		return errors.New("fieldsMap is nil")
	}

	if t.fields == nil {
		return errors.New("fields is nil")
	}

	// First field?
	if len(t.fields) == 0 {
		_, ok := f.(*FieldMetaUint64)
		if !ok {
			return errors.New("First fieldmeta added must be primary key type FieldMetaUint64")
		}
	}

	t.fieldsMap[f.Name()] = f
	t.fields = append(t.fields, f)

	return nil
}

func (t *TableDef) createTableString(dialect Dialect) (string, error) {
	if !t.done {
		return "", errors.New("Table must be done")
	}
	if dialect == nil {
		return "", errors.New("Dialect is nil")
	}
	return dialect.CreateTableString(t)
}

func (t *TableDef) CreatePreparedStatementInsertAllFields(dialect Dialect) (string, error) {
	a, b := t.CreatePreparedStatementInsertSomeFields(dialect, t.fields...)
	return a, b
}

func (t *TableDef) CreatePreparedStatementInsertFromRecord(dialect Dialect, record *Record) (string, error) {
	if record == nil {
		return "", errors.New("Record cannot be nil")
	}
	fields := make([]FieldMeta, 0)
	for i, _ := range record.values {
		fields = append(fields, record.tableDef.fields[i])
	}
	err := t.SetDone()
	if err != nil {
		return "", err
	}

	return t.CreatePreparedStatementInsertSomeFields(dialect, fields...)
}

func (t *TableDef) CreatePreparedStatementInsertSomeFields(dialect Dialect, fields ...FieldMeta) (string, error) {
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
		preparedValueFormat, err := preparedValueFormat(dialect, i)
		if err != nil {
			return "", err
		}
		values += preparedValueFormat
	}

	st += ")"
	values += ")"

	st = st + " VALUES " + values
	return st, nil
}
