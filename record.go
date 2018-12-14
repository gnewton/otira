package otira

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
)

type Record struct {
	values     []interface{}
	valueIsSet []bool
	tableMeta  *TableMeta
	fields     []FieldMeta
	fieldsMap  map[string]int
	//fieldsMap       map[string]FieldMeta
	validating      bool
	tx              *sql.Tx
	stmt            *sql.Stmt
	preparedString  string
	relationRecords []*RelationRecord
}

type RelationRecord struct {
	record   *Record
	relation Relation
}

func (r *Record) Prepare(tx *sql.Tx) error {

	if r.tx == tx {
		if r.stmt != nil {
			return nil
		}
	} else {
		r.tx = tx
	}

	var err error
	if r.preparedString == "" {
		r.preparedString, err = r.tableMeta.CreatePreparedStatementInsertFromRecord(new(DialectSqlite3), r)
		if err != nil {
			return err
		}
	}
	r.stmt, err = r.tx.Prepare(r.preparedString)
	if err != nil {
		return err
	}
	return nil

}

func (r *Record) Values() []interface{} {
	return r.values
}

func newRecord(tm *TableMeta, fields []FieldMeta, stmt *sql.Stmt) (*Record, error) {
	if tm == nil {
		return nil, errors.New("TableMeta is nil")
	}
	if fields == nil {
		return nil, errors.New("Fields is nil")
	}
	if len(fields) == 0 {
		return nil, errors.New("Fields is zero length")
	}
	rec := new(Record)
	rec.tableMeta = tm

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
	log.Println("AddRelationRecord")
	log.Println(rel)
	log.Println(record)
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
	switch v := rel.(type) {
	case *OneToMany:
		log.Println("===== " + v.String())

	case *ManyToMany:
		log.Println(v.String())
	}
	log.Println(r.relationRecords[0].record)
	log.Println(r.relationRecords[0].relation)
	log.Println("END AddRelationRecord")
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
	return newRecord(r.tableMeta, r.fields, r.stmt)
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
	if r.validating {
		if !r.fields[i].IsSameType(v) {
			return errors.New("Incorrect type")
		}
	}
	r.values[i] = v
	r.valueIsSet[i] = true
	return nil
}

func (r *Record) String() string {
	var s string

	s += "TableName:" + r.tableMeta.GetName()
	for i := 0; i < len(r.fields); i++ {
		s += "\n " + r.fields[i].Name() + ":" + toString(r.values[i])
	}
	return s
}

//func (r *Record) SetPrimaryKey() error {

//}
