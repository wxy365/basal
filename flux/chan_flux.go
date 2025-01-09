package flux

import (
	"github.com/wxy365/basal/fn"
	"github.com/wxy365/basal/iterator"
	"github.com/wxy365/basal/opt"
)

type chanFlux[T any] struct {
	data chan T
}

func (f *chanFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *chanFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *chanFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *chanFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumedFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *chanFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *chanFlux[T]) Skip(n uint) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *chanFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *chanFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		next := itr.Next()
		consumer(next)
	}
}

func (f *chanFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *chanFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *chanFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *chanFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *chanFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *chanFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *chanFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *chanFlux[T]) Emit() opt.Opt[T] {
	d, ok := <-f.data
	if ok {
		return opt.Of(d)
	}
	return opt.Empty[T]()
}

func (f *chanFlux[T]) Iterator() iterator.Iterator[T] {
	return iterator.FromChan(f.data)
}

func (f *chanFlux[T]) Close() {
	close(f.data)
}
