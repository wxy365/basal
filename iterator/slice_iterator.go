package iterator

import "github.com/wxy365/basal/fn"

// sliceIterator is the iterator for transversing slice
// non thread safe
type sliceIterator[T any] struct {
	tgt       []T
	cur       int
	removable bool
}

func (s *sliceIterator[T]) HasNext() bool {
	return s.cur < len(s.tgt)
}

func (s *sliceIterator[T]) Next() T {
	r := s.tgt[s.cur]
	s.cur++
	s.removable = true
	return r
}

func (s *sliceIterator[T]) Remove() {
	if s.removable {
		s.cur--
		s.tgt = append(s.tgt[:s.cur], s.tgt[s.cur+1:]...)
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
