package otira

import (
	"errors"
)

type Field struct {
	field       interface{}
	fieldDef    FieldDef
	hasSetValue bool
}

func (f *Field) SetValueFast(v interface{}) {
	f.field = v
	f.hasSetValue = true
}

func (f *Field) SetValue(v interface{}) error {
	f.field = v
	f.hasSetValue = true
	if !f.fieldDef.IsSameType(v) {
		return errors.New("Incorrect type:" + toString(v))
	}
	return nil
}
