package flux

import (
	"github.com/wxy365/basal/fn"
	"github.com/wxy365/basal/iterator"
	"github.com/wxy365/basal/opt"
)

type sliceFlux[T any] struct {
	tgt  []T
	next int
}

func (s *sliceFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: s,
		filter: predicate,
	}
}

func (s *sliceFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     s,
		comparator: comparator,
	}
}

func (s *sliceFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   s,
		comparer: comparer,
	}
}

func (s *sliceFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumedFlux[T]{
		origin:   s,
		consumer: consumer,
	}
}

func (s *sliceFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  s,
		maxSize: size,
	}
}

func (s *sliceFlux[T]) Skip(n uint) Flux[T] {
	return &skippedFlux[T]{
		origin: s,
		skip:   n,
	}
}

func (s *sliceFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    s,
		closeHook: hook,
	}
}

func (s *sliceFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := s.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (s *sliceFlux[T]) ToSlice() []T {
	return s.tgt
}

func (s *sliceFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](s, comparator)
}

func (s *sliceFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](s, comparator)
}

func (s *sliceFlux[T]) Count() int64 {
	return count[T](s)
}

func (s *sliceFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](s, predicate)
}

func (s *sliceFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](s, predicate)
}

func (s *sliceFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](s, predicate)
}

func (s *sliceFlux[T]) Emit() opt.Opt[T] {
	if s.next > len(s.tgt)-1 {
		return opt.Empty[T]()
	}
	ret := opt.Of(s.tgt[s.next])
	s.next++
	return ret
}

func (s *sliceFlux[T]) Iterator() iterator.Iterator[T] {
	if s.next > len(s.tgt)-1 {
		return iterator.OfSlice([]T{})
	}
	return iterator.OfSlice(s.tgt[s.next:])
}

func (s *sliceFlux[T]) Close() {

}
