package fn

import "sync/atomic"

type Counter int64

func (c *Counter) Incr() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *Counter) GetAndIncr() int64 {
	ret := *c
	atomic.AddInt64((*int64)(c), 1)
	return int64(ret)
}

func (c *Counter) Set(val int64) {
	atomic.StoreInt64((*int64)(c), val)
}
