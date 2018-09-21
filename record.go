package otira

import (
	"errors"
	"strconv"
)

type Record struct {
	values     []interface{}
	tableMeta  *TableMeta
	fields     []FieldMeta
	fieldsMap  map[string]int
	validating bool
}

func newRecord(tm *TableMeta, fields []FieldMeta) (*Record, error) {
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
	rec.fields = fields
	rec.fieldsMap = make(map[string]int, len(rec.fields))
	for i := 0; i < len(fields); i++ {
		rec.fieldsMap[fields[i].Name()] = i
	}
	return rec, nil
}

func (r *Record) Clone() (*Record, error) {
	return newRecord(r.tableMeta, r.fields)
}

func (r *Record) SetByName(f string, v interface{}) error {
	i, ok := r.fieldsMap[f]
	if !ok {
		return errors.New("Field with name " + f + " does not exist")
	}
	r.values[i] = v
	return nil
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

	return nil
}
