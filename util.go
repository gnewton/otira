package otira

import (
	"errors"
)

func baseFieldMetaErrors(f FieldMeta) error {
	if f == nil {
		return errors.New("Field cannot be nil")
	}
	if f.Name() == "" {
		return errors.New("Field cannot have nil name")
	}
	return nil
}
