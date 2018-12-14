package otira

import (
	"errors"
	"log"
)

type joinCache struct {
	joinKeys map[string]uint64
}

func NewJoinCache() *joinCache {
	jc := new(joinCache)
	jc.joinKeys = make(map[string]uint64)
	return jc
}

func (jc *joinCache) GetJoinKey(r *Record) (uint64, bool, error) {
	cacheKey, err := makeKey(r)
	if err != nil {
		return 0, true, err
	}
	joinKey, exists := jc.joinKeys[cacheKey]
	if !exists {
		// New record: Get the next primary key for the table
		if r.tableMeta.useRecordPrimaryKeys {
			var ok bool
			index, ok := r.fieldsMap[r.tableMeta.primaryKey.Name()]
			log.Println("Primary key field name=" + r.tableMeta.primaryKey.Name())
			log.Println(index)
			if !ok {
				return 0, true, errors.New("Primary key:" + r.tableMeta.primaryKey.Name() + " is not in fieldmeta map")
			}
			switch k := r.values[index].(type) {
			case uint64:
				joinKey = k
			default:
				return 0, true, errors.New("Primary key:" + r.tableMeta.primaryKey.Name() + " is not type uint64" + "  " + toString(k))
			}
		} else {
			joinKey, err = r.tableMeta.Next()
			if err != nil {
				return 0, true, err
			}
			r.values[r.fieldsMap[r.tableMeta.primaryKey.Name()]] = joinKey
		}
		jc.joinKeys[cacheKey] = joinKey
	}

	return joinKey, exists, nil

}

func (jc *joinCache) getJoinKey(r *Record) {
	_, _ = makeKey(r)
}

// Makes a key string from concatenating the discrimFields' values
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
