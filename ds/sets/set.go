package sets

import (
	"github.com/wxy365/basal/fn"
	"github.com/wxy365/basal/iterator"
)

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(t T) {
	s.m[t] = struct{}{}
}

func (s *Set[T]) Remove(t T) {
	delete(s.m, t)
}

func (s *Set[T]) Length() int {
	return len(s.m)
}

func (s *Set[T]) Contains(t T) bool {
	_, exists := s.m[t]
	return exists
}

func (s *Set[T]) Iterator() iterator.Iterator[T] {
	if s.m == nil || s.Length() == 0 {
		return &iterator.DummyIterator[T]{}
	}
	pipe := make(chan T, s.Length())
	go func() {
		defer func() {
			// keep silence
			_ = recover()
		}()
		for k := range s.m {
			pipe <- k
		}
		close(pipe)
	}()
	itr := &setIterator[T]{
		tgt:  s,
		pipe: pipe,
	}
	itr.next, itr.hasNext = <-pipe
	return itr
}

type setIterator[T comparable] struct {
	tgt     *Set[T]
	pipe    chan T
	cur     T
	next    T
	hasNext bool
}

func (s *setIterator[T]) HasNext() bool {
	return s.hasNext
}

func (s *setIterator[T]) Next() T {
	s.cur = s.next
	s.next, s.hasNext = <-s.pipe
	return s.cur
}

func (s *setIterator[T]) Remove() {
	s.tgt.Remove(s.cur)
}

func (s *setIterator[T]) ForEach(consumer fn.Consumer[T]) {
	for s.HasNext() {
		consumer(s.Next())
	}
}

func (s *setIterator[T]) Close() {
	close(s.pipe)
}
