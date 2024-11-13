package iterator

import "github.com/wxy365/basal/fn"

func Map[T, R any](itr Iterator[T], mapper fn.Function[T, R]) Iterator[R] {
	return &mappedIterator[T, R]{
		origin: itr,
		mapper: mapper,
	}
}

type mappedIterator[T, R any] struct {
	origin Iterator[T]
	mapper fn.Function[T, R]
}

func (m *mappedIterator[T, R]) HasNext() bool {
	return m.origin.HasNext()
}

func (m *mappedIterator[T, R]) Next() R {
	return m.mapper(m.origin.Next())
}

func (m *mappedIterator[T, R]) Remove() {
	m.origin.Remove()
}

func (m *mappedIterator[T, R]) ForEach(consumer fn.Consumer[R]) {
	for m.HasNext() {
		consumer(m.Next())
	}
}
