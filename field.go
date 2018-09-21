package otira

import (
	"errors"
)

type Field struct {
	field       interface{}
	fieldMeta   FieldMeta
	hasSetValue bool
}

func (f *Field) SetValueFast(v interface{}) {
	f.field = v
	f.hasSetValue = true
}

func (f *Field) SetValue(v interface{}) error {
	f.field = v
	f.hasSetValue = true
	if !f.fieldMeta.IsSameType(v) {
		return errors.New("Incorrect type")
	}
	return nil
}
