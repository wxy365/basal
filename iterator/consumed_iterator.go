package iterator

import "github.com/wxy365/basal/fn"

func Consume[T any](itr Iterator[T], consumer fn.Consumer[T]) Iterator[T] {
	return &consumedIterator[T]{
		origin:   itr,
		consumer: consumer,
	}
}

type consumedIterator[T any] struct {
	origin   Iterator[T]
	consumer fn.Consumer[T]
}

func (c *consumedIterator[T]) HasNext() bool {
	return c.origin.HasNext()
}

func (c *consumedIterator[T]) Next() T {
	n := c.origin.Next()
	c.consumer(n)
	return n
}

func (c *consumedIterator[T]) Remove() {
	c.origin.Remove()
}

func (c *consumedIterator[T]) ForEach(consumer fn.Consumer[T]) {
	for c.HasNext() {
		consumer(c.Next())
	}
}
