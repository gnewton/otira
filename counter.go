package otira

import (
	"log"
	"sync/atomic"
)

type Counter interface {
	Next() (uint64, error)
	Value() (uint64, error)
}

type ICounter struct {
	counter uint64
}

func (c *ICounter) Next() (uint64, error) {
	log.Println("Counter: " + toString(c.counter))
	atomic.AddUint64(&c.counter, 1)
	return c.counter, nil
}
