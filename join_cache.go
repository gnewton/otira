package otira

import (
	"errors"
)

type joinCache struct {
	keys map[string]int64
}

func (jc *joinCache) getJoinKey(r *Record) {

	_, _ = makeKey(r)
	//return nil
}

func makeKey(r *Record) (string, error) {
	var keyString string

	flen := len(r.fields)
	dflen := len(r.tableMeta.discrimFields)
	for i, _ := range r.tableMeta.discrimFields {
		fm := r.tableMeta.discrimFields[i]
		j, ok := r.fieldsMap[fm.Name()]
		if !ok {
			return "", errors.New("Field name [" + fm.Name() + "] is not a field in table " + r.tableMeta.GetName())
		}
		if i < 0 || i > flen {
			return "", errors.New("Index out of bounds for field [" + fm.Name() + "]")
		}
		if r.values[j] == nil {
			return "", errors.New("toString failed for value:" + toString(j))
		}
		if dflen == 1 {
			keyString = toString(r.values[j])
		} else {
			keyString += "_" + toString(r.values[j]) + "|"
		}
	}
	return keyString, nil

}
