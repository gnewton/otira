package otira

import (
	"errors"
	"strconv"
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

func findRelationPK(record *Record, rel *OneToMany) (string, error) {
	var k string
	for i := 0; i < len(rel.rightTableUniqueFields); i++ {
		fname := rel.rightTableUniqueFields[i].Name()
		lookup, ok := record.fieldsMap[fname]

		if !ok {
			return "", errors.New("Field " + fname + " not in relation " + rel.Name())
		}
		if lookup >= len(record.values) {
			return "", errors.New("Index too large " + strconv.Itoa(lookup) + "; should be <" + strconv.Itoa(len(record.values)))
		}
		fvalue := toString(record.values[lookup])
		k = k + fname + "_" + fvalue + " |"
	}

	return k, nil
}

func toString(t interface{}) string {
	switch v := t.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	}
	return "ERROR"
}
