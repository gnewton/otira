package otira

import (
	"sync/atomic"
)

// type Counter interface {
// 	Next() (int64, error)
// 	Value() (int64, error)
// }

type ICounter struct {
	counter int64
}

func (c *ICounter) Next() int64 {
	atomic.AddInt64(&c.counter, 1)
	return c.counter
}
