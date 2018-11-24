package otira

import (
	"sync/atomic"
)

type Counter interface {
	Next() (uint64, error)
}

type ICounter struct {
	value uint64
}

func (c *ICounter) Next() (uint64, error) {
	atomic.AddUint64(&c.value, 1)
	return c.value, nil
}
