package otira

import (
	"encoding/binary"
)

type Set interface {
	Put(key int64) error
	Contains(key int64) (bool, error)
	Close() error
}

func int64ToByteArray(key int64) []byte {
	bk := make([]byte, 8)
	binary.LittleEndian.PutUint64(bk, uint64(key))

	return bk
}
