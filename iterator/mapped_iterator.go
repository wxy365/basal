package iterator

import "github.com/wxy365/basal/fn"

func Map[R, T any](itr Iterator[R], mapper fn.Function[R, T]) Iterator[T] {
	return &mappedIterator[R, T]{
		origin: itr,
		mapper: mapper,
	}
}

type mappedIterator[R, T any] struct {
	origin Iterator[R]
	mapper fn.Function[R, T]
}

func (m *mappedIterator[R, T]) HasNext() bool {
	return m.origin.HasNext()
}

func (m *mappedIterator[R, T]) Next() T {
	return m.mapper(m.origin.Next())
}

func (m *mappedIterator[R, T]) Remove() {
	m.origin.Remove()
}

func (m *mappedIterator[R, T]) ForEach(consumer fn.Consumer[T]) {
	for m.HasNext() {
		consumer(m.Next())
	}
}
