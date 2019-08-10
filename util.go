package otira

import (
	"errors"
	"strconv"
)

func baseFieldDefErrors(f FieldDef) error {
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
	for i := 0; i < len(rel.RightTableUniqueFields); i++ {
		fname := rel.RightTableUniqueFields[i].Name()
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

func supportedType(t interface{}) bool {
	switch t.(type) {
	case string, int64, uint32, int, int8, int16, int32, uint64, float32, float64, []byte, bool:
		return true
	}
	return false

}

const STRING_ERROR = " -ERROR- "

func toString(t interface{}) string {
	if t == nil {
		return " --nil-- "
	}

	switch v := t.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(int64(v), 10)

	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)

	case uint:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)

	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)

	case bool:
		return strconv.FormatBool(v)

	}
	return STRING_ERROR
}

func unsetPrimaryKey(rec *Record) bool {
	return !rec.valueIsSet[0]
}
