package otira

import (
	"errors"
)

type Index struct {
	name      string
	fieldDefs []*FieldDef
	Unique    bool
}

func NewIndex(name string, fields ...*FieldDef) (*Index, error) {
	if name == "" {
		return nil, errors.New("String cannot be empty")
	}
	if fields == nil {
		return nil, errors.New("Fields cannot be nil")
	}

	if len(fields) == 0 {
		return nil, errors.New("Fields cannot be length 0")
	}

	for i := 0; i < len(fields); i++ {
		if fields[i] == nil {
			return nil, errors.New("Field[" + toString(i) + "] is nil")
		}
	}

	index := new(Index)
	index.name = name
	index.fieldDefs = fields

	return index, nil

}
