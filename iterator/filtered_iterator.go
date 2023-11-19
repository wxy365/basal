package iterator

import "github.com/wxy365/basal/fn"

func Filter[T any](itr Iterator[T], predicate fn.Predicate[T]) Iterator[T] {
	f := &filteredIterator[T]{
		origin:    itr,
		predicate: predicate,
	}
	f.stepNext()
	return f
}

type filteredIterator[T any] struct {
	origin    Iterator[T]
	predicate fn.Predicate[T]
	next      *T
}

func (f *filteredIterator[T]) HasNext() bool {
	return f.next != nil
}

func (f *filteredIterator[T]) Next() T {
	next := f.next
	f.next = nil
	f.stepNext()
	return *next
}

func (f *filteredIterator[T]) stepNext() {
	for f.origin.HasNext() {
		newNext := f.origin.Next()
		if f.predicate(newNext) {
			f.next = &newNext
			break
		}
	}
}

func (f *filteredIterator[T]) Remove() {
	panic("cannot remove element over a filtered iterator")
}

func (f *filteredIterator[T]) ForEach(consumer fn.Consumer[T]) {
	for f.HasNext() {
		consumer(f.Next())
	}
}
