package flux

import "github.com/wxy365/basal/fn"

func newFluxIterator[T any](pipe chan T) *fluxIterator[T] {
	itr := &fluxIterator[T]{
		pipe: pipe,
	}
	if n, ok := <-pipe; ok {
		itr.next = &n
	}
	return itr
}

type fluxIterator[T any] struct {
	pipe chan T
	next *T
}

func (f *fluxIterator[T]) HasNext() bool {
	return f.next != nil
}

func (f *fluxIterator[T]) Next() T {
	var t T
	if f.next == nil {
		return t
	}
	next := *f.next
	if n, ok := <-f.pipe; ok {
		f.next = &n
	} else {
		f.next = nil
	}
	return next
}

func (f *fluxIterator[T]) Remove() {

}

func (f *fluxIterator[T]) ForEach(consumer fn.Consumer[T]) {
	for f.HasNext() {
		consumer(f.Next())
	}
}
