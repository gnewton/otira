package otira

import (
	"errors"
	"log"
)

type Field struct {
	field       interface{}
	fieldMeta   FieldMeta
	hasSetValue bool
}

func (f *Field) SetValue(v interface{}) error {
	f.field = v
	f.hasSetValue = true
	log.Println(f.fieldMeta.IsSameType(v))
	if !f.fieldMeta.IsSameType(v) {
		return errors.New("Incorrect type")
	}
	return nil
}
