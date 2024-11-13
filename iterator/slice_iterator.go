package iterator

import "github.com/wxy365/basal/fn"

// sliceIterator is the iterator for transversing slice
// non thread safe
type sliceIterator[T any] struct {
	tgt       *[]T
	next      int
	removable bool
}

func (s *sliceIterator[T]) HasNext() bool {
	return s.next < len(*s.tgt)
}

func (s *sliceIterator[T]) Next() T {
	if s.next >= len(*s.tgt) {
		panic("calling method 'Next' on a finished iterator")
	}
	r := (*s.tgt)[s.next]
	s.next++
	s.removable = true
	return r
}

func (s *sliceIterator[T]) Remove() {
	if s.removable {
		s.next--
		*s.tgt = append((*s.tgt)[:s.next], (*s.tgt)[s.next+1:]...)
		s.removable = false
	}
}

func (s *sliceIterator[T]) removeLast() {
	if s.removable {
		s.next--
		*s.tgt = append((*s.tgt)[:s.next], (*s.tgt)[s.next+1:]...)
		s.removable = false
	}
}

func (s *sliceIterator[T]) ForEach(consumer fn.Consumer[T]) {
	for s.HasNext() {
		consumer(s.Next())
	}
}

func (s *sliceIterator[T]) Close() {
}
