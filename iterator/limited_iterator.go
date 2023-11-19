package iterator

import "github.com/wxy365/basal/fn"

func Limit[T any](itr Iterator[T], maxSize uint) Iterator[T] {
	return &limitedIterator[T]{
		origin:  itr,
		maxSize: maxSize,
	}
}

type limitedIterator[T any] struct {
	origin  Iterator[T]
	maxSize uint
	count   uint
}

func (l *limitedIterator[T]) HasNext() bool {
	return l.origin.HasNext() && l.count < l.maxSize
}

func (l *limitedIterator[T]) Next() T {
	if l.count >= l.maxSize {
		panic("calling method 'Next' on a finished iterator")
	}
	next := l.origin.Next()
	l.count++
	return next
}

func (l *limitedIterator[T]) Remove() {
	l.origin.Remove()
}

func (l *limitedIterator[T]) ForEach(consumer fn.Consumer[T]) {
	for l.HasNext() {
		consumer(l.Next())
	}
}
