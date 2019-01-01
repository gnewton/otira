package otira

import (
	"encoding/binary"
)

type Set interface {
	Put(key uint64) error
	Contains(key uint64) (bool, error)
	Close() error
}

func uint64ToByteArray(key uint64) []byte {
	bk := make([]byte, 8)
	binary.LittleEndian.PutUint64(bk, key)

	return bk
}
