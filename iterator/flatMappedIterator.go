package iterator

import "github.com/wxy365/basal/fn"

func FlatMap[R, T any](itr Iterator[R], mapper fn.Function[R, Iterator[T]]) Iterator[T] {
	return &flatMappedIterator[R, T]{
		origin: itr,
		mapper: mapper,
	}
}

type flatMappedIterator[R, T any] struct {
	origin    Iterator[R]
	mapper    fn.Function[R, Iterator[T]]
	mappedEle Iterator[T]
}

func (f *flatMappedIterator[R, T]) HasNext() bool {
	if f.mappedEle == nil {
		if !f.origin.HasNext() {
			return false
		}
		r := f.origin.Next()
		f.mappedEle = f.mapper(r)
	} else {
		if f.mappedEle.HasNext() {
			return true
		}
		f.mappedEle = nil
	}
	return f.HasNext()
}

func (f *flatMappedIterator[R, T]) Next() T {
	if !f.HasNext() {
		panic("calling 'Next' method on finished iterator")
	}
	if f.mappedEle == nil {
		f.mappedEle = f.mapper(f.origin.Next())
	}
	return f.mappedEle.Next()
}

func (f *flatMappedIterator[R, T]) Remove() {
	f.mappedEle.Remove()
}

func (f *flatMappedIterator[R, T]) ForEach(consumer fn.Consumer[T]) {
	for f.HasNext() {
		consumer(f.Next())
	}
}
