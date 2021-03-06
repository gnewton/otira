package otira

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
)

type saveStatus int

const (
	unsaved saveStatus = iota
	inPreparation
	saved
)

const UnassignedPK = -99

type Record struct {
	//fieldsMap       map[string]FieldDef
	fields          []FieldDef
	fieldsMap       map[string]int
	preparedString  string
	relationRecords []*RelationRecord
	stmt            *sql.Stmt
	tableDef        *TableDef
	tx              *sql.Tx
	valueIsSet      []bool
	values          []interface{}
	status          saveStatus
}

type RelationRecord struct {
	record   *Record
	relation Relation
}

func (r *Record) Values() []interface{} {
	return r.values
}

func newRecord(tm *TableDef, fields []FieldDef, stmt *sql.Stmt) (*Record, error) {
	if tm == nil {
		return nil, errors.New("TableDef is nil")
	}
	if fields == nil {
		return nil, errors.New("Fields is nil")
	}
	if len(fields) == 0 {
		return nil, errors.New("Fields is zero length")
	}
	rec := new(Record)
	rec.tableDef = tm

	rec.values = make([]interface{}, len(fields))
	rec.valueIsSet = make([]bool, len(fields))
	rec.fields = fields
	if stmt == nil {
		rec.stmt = stmt
	} else {
		tm.NewRecordSomeFields(fields...)
	}

	//rec.fieldsMap = make(map[string]int, len(rec.fields))
	rec.fieldsMap = make(map[string]int, len(rec.fields))
	for i := 0; i < len(fields); i++ {
		rec.fieldsMap[fields[i].Name()] = i
	}
	return rec, nil
}

func (r *Record) AddRelationRecord(rel Relation, record *Record) error {
	if rel == nil {
		return errors.New("Relation is nil")
	}

	if record == nil {
		return errors.New("Record is nil")
	}

	relationRecord := new(RelationRecord)
	relationRecord.record = record
	relationRecord.relation = rel
	r.relationRecords = append(r.relationRecords, relationRecord)
	// TODO

	// switch v := rel.(type) {
	// case *OneToMany:
	// 	log.Println("===== " + v.String())

	// case *ManyToMany:
	// 	log.Println(v.String())
	// }
	return nil
}

func (r *Record) Reset() error {
	if r.values == nil {
		return errors.New("Values is nil")
	}
	for i := 0; i < len(r.values); i++ {
		r.values[i] = nil
	}
	return nil
}

func (r *Record) Clone() (*Record, error) {
	return newRecord(r.tableDef, r.fields, r.stmt)
}

func (r *Record) SetByName(f string, v interface{}) error {
	i, ok := r.fieldsMap[f]
	if !ok {
		return errors.New("Field with name " + f + " does not exist")
	}

	if !supportedType(v) {
		return errors.New("Added value type for field [" + f + "] is not supported")
	}

	return r.Set(i, v)
}

func (r *Record) Set(i int, v interface{}) error {

	if i < 0 || i > len(r.values) {
		return errors.New("Index out of bounds. Should be 0.." + strconv.Itoa(len(r.values)) + "; Actual: " + strconv.Itoa(i))
	}
	if r.tableDef.validating {
		if !r.fields[i].IsSameType(v) {
			return errors.New("Incorrect type")
		}
	}
	r.values[i] = v
	r.valueIsSet[i] = true
	log.Println("true")
	return nil
}

func (r *Record) PrimaryKeyValue() (int64, error) {
	switch v := r.values[0].(type) {
	case int64:
		return v, nil
	}
	log.Println("+++++++++++")
	log.Println(r.values[0])
	log.Println("+++++++++++")
	return 0, errors.New("Primary key is not int64")
}

func (r *Record) String() string {
	var s string

	s += "TableName:" + r.tableDef.Name()
	for i := 0; i < len(r.fields); i++ {
		s += "\n " + r.fields[i].Name() + ":" + toString(r.values[i])
	}
	return s
}

//func (r *Record) SetPrimaryKey() error {

//}
