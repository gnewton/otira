package otira

import (
	"errors"
	//	"log"
)

type joinCache struct {
	joinKeys map[string]int64
}

func NewJoinCache() *joinCache {
	jc := new(joinCache)
	jc.joinKeys = make(map[string]int64)
	return jc
}

func (jc *joinCache) MakeJoinKey(r *Record) (int64, bool, error) {
	cacheKey, err := makeKey(r)
	if err != nil {
		return 0, true, err
	}
	joinKey, exists := jc.joinKeys[cacheKey]
	if !exists {
		// New record: Get the next primary key for the table
		if r.tableDef.UseRecordPrimaryKeys {
			pk, ok := r.values[0].(int64)
			if !ok {
				return 0, false, errors.New("Primary key value is not int64; table=" + r.tableDef.Name())
			} else {
				joinKey = pk
			}

		} else {
			joinKey = r.tableDef.Next()
			r.values[0] = joinKey
		}
		jc.joinKeys[cacheKey] = joinKey
	} else {
		//log.Println("---> Cache hit:[" + cacheKey + "]")
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
	dflen := len(r.tableDef.joinDiscrimFields)

	if dflen == 0 {
		return "", errors.New("No discrim fields for table:" + r.tableDef.name)

	}

	for i, _ := range r.tableDef.joinDiscrimFields {
		fm := r.tableDef.joinDiscrimFields[i]
		j, ok := r.fieldsMap[fm.Name()]
		if !ok {
			return "", errors.New("Field name [" + fm.Name() + "] is not a field in table " + r.tableDef.Name())
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
