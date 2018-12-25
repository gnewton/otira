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
			pk, ok := r.values[0].(uint64)
			if !ok {
				return 0, false, errors.New("Primary key value is not uint64; table=" + r.tableMeta.GetName())
			} else {
				joinKey = pk
			}

		} else {
			joinKey, err = r.tableMeta.Next()
			if err != nil {
				return 0, true, err
			}
			r.values[0] = joinKey
		}
		jc.joinKeys[cacheKey] = joinKey
	} else {
		log.Println("---> Cache hit:[" + cacheKey + "]")
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

	if dflen == 0 {
		return "", errors.New("No discrim fields for table:" + r.tableMeta.name)

	}

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
	if keyString == "" {
		return "", errors.New("key cache is empty string")
	}

	return keyString, nil

}
