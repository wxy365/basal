package iterator

import "github.com/wxy365/basal/fn"

func FlatMap[T, R any](itr Iterator[T], mapper fn.Function[T, Iterator[R]]) Iterator[R] {
	return &flatMappedIterator[T, R]{
		origin: itr,
		mapper: mapper,
	}
}

type flatMappedIterator[T, R any] struct {
	origin    Iterator[T]
	mapper    fn.Function[T, Iterator[R]]
	mappedEle Iterator[R]
}

func (f *flatMappedIterator[T, R]) HasNext() bool {
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

func (f *flatMappedIterator[T, R]) Next() R {
	if !f.HasNext() {
		panic("calling 'Next' method on finished iterator")
	}
	if f.mappedEle == nil {
		f.mappedEle = f.mapper(f.origin.Next())
	}
	return f.mappedEle.Next()
}

func (f *flatMappedIterator[T, R]) Remove() {
	f.mappedEle.Remove()
}

func (f *flatMappedIterator[T, R]) ForEach(consumer fn.Consumer[R]) {
	for f.HasNext() {
		consumer(f.Next())
	}
}
