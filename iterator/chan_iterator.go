package iterator

import "github.com/wxy365/basal/fn"

type chanIterator[T any] struct {
	ch   chan T
	next *T
}

func (c *chanIterator[T]) HasNext() bool {
	return c.next != nil
}

func (c *chanIterator[T]) Next() T {
	var t T
	if c.next == nil {
		return t
	}
	next := *c.next
	if n, ok := <-c.ch; ok {
		c.next = &n
	} else {
		c.next = nil
	}
	return next
}

func (c *chanIterator[T]) Remove() {
	// Once an element is read from the channel, it is also removed from the channel
}

func (c *chanIterator[T]) ForEach(consumer fn.Consumer[T]) {
	for c.HasNext() {
		consumer(c.Next())
	}
}
