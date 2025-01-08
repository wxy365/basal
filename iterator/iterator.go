package iterator

import (
	"github.com/wxy365/basal/fn"
)

type Iterator[T any] interface {
	HasNext() bool
	Next() T
	Remove() // delete the last accessed item from the underlying data structure, maybe a lists or a map
	ForEach(consumer fn.Consumer[T])
}

type DummyIterator[T any] struct {
}

func (d *DummyIterator[T]) HasNext() bool {
	return false
}

func (d *DummyIterator[T]) Next() T {
	var t T
	return t
}

func (d *DummyIterator[T]) Remove() {
}

func (d *DummyIterator[T]) ForEach(consumer fn.Consumer[T]) {
}

func OfMap[K comparable, V any](tgt map[K]V) Iterator[MapEntry[K, V]] {
	if len(tgt) == 0 {
		return &DummyIterator[MapEntry[K, V]]{}
	}
	pipe := make(chan MapEntry[K, V], len(tgt))
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
	itr := &mapIterator[K, V]{
		tgt:  tgt,
		pipe: pipe,
	}
	itr.next, itr.hasNext = <-pipe
	return itr
}

func OfSlice[T any](tgt []T) Iterator[T] {
	if len(tgt) == 0 {
		return &DummyIterator[T]{}
	}
	return &sliceIterator[T]{
		tgt:       &tgt,
		next:      0,
		removable: false,
	}
}

func OfChan[T any](tgt chan T) Iterator[T] {
	if tgt == nil {
		return &DummyIterator[T]{}
	}
	next := <-tgt
	return &chanIterator[T]{
		ch:   tgt,
		next: &next,
	}
}
