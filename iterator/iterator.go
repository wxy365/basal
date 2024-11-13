package iterator

import "github.com/wxy365/basal/fn"

type Iterator[T any] interface {
	HasNext() bool
	Next() T
	Remove() // delete the last accessed item from the underlying data structure, maybe a lists or a map
	ForEach(consumer fn.Consumer[T])
}

type dummyIterator[T any] struct {
}

func (d *dummyIterator[T]) HasNext() bool {
	return false
}

func (d *dummyIterator[T]) Next() T {
	var t T
	return t
}

func (d *dummyIterator[T]) Remove() {
}

func (d *dummyIterator[T]) ForEach(consumer fn.Consumer[T]) {
}

func OfMap[K comparable, V any](tgt map[K]V) Iterator[MapEntry[K, V]] {
	if tgt == nil {
		return &dummyIterator[MapEntry[K, V]]{}
	}
	pipe := make(chan MapEntry[K, V])
	go func() {
		defer func() {
			// keep silence
			_ = recover()
		}()
		for k, v := range tgt {
			pipe <- &defaultMapEntry[K, V]{key: k, value: v}
		}
		close(pipe)
	}()
	return &mapIterator[K, V]{
		tgt:  tgt,
		pipe: pipe,
		next: <-pipe,
	}
}

func OfSlice[T any](tgt []T) Iterator[T] {
	if tgt == nil {
		return &dummyIterator[T]{}
	}
	return &sliceIterator[T]{
		tgt:       &tgt,
		next:      0,
		removable: false,
	}
}

func OfChan[T any](tgt chan T) Iterator[T] {
	if tgt == nil {
		return &dummyIterator[T]{}
	}
	next := <-tgt
	return &chanIterator[T]{
		ch:   tgt,
		next: &next,
	}
}
