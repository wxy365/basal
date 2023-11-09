package flux

import (
	"github.com/wxy365/basal/fn"
	"github.com/wxy365/basal/iterator"
	"github.com/wxy365/basal/opt"
)

type Flux[T any] interface {
	Filter(predicate fn.Predicate[T]) Flux[T]
	Sorted(comparator fn.Comparator[T]) Flux[T]
	Distinct(comparer fn.Comparer[T]) Flux[T]
	OnEmit(consumer fn.Consumer[T]) Flux[T]
	Limit(size int64) Flux[T]
	Skip(n int64) Flux[T]
	OnClose(func()) Flux[T]
	ForEach(consumer fn.Consumer[T])
	ToSlice() []T
	Min(comparator fn.Comparator[T]) opt.Opt[T]
	Max(comparator fn.Comparator[T]) opt.Opt[T]
	Count() int64
	AnyMatch(predicate fn.Predicate[T]) bool
	AllMatch(predicate fn.Predicate[T]) bool
	NoneMatch(predicate fn.Predicate[T]) bool
	Emit() opt.Opt[T]
	Iterator() iterator.Iterator[T]
	Close()
}

func Map[R, T any](f Flux[R], mapper fn.Function[R, T]) Flux[T] {
	return &mappedFlux[R, T]{
		origin: f,
		mapper: mapper,
	}
}

func FlatMap[R, T any](f Flux[R], mapper fn.Function[R, Flux[T]]) Flux[T] {
	return &flatMappedFlux[R, T]{
		origin: f,
		mapper: mapper,
	}
}

func Collect[T, A, R any](f Flux[T], collector Collector[T, A, R]) R {
	container := collector.Supplier()()
	f.ForEach(func(t T) {
		collector.Gatherer()(container, t)
	})
	return collector.Finisher()(container)
}

type fluxImpl[T any] struct {
	data chan T
}

func (f *fluxImpl[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *fluxImpl[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *fluxImpl[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *fluxImpl[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *fluxImpl[T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *fluxImpl[T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *fluxImpl[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *fluxImpl[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		next := itr.Next()
		consumer(next)
	}
}

func (f *fluxImpl[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *fluxImpl[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *fluxImpl[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *fluxImpl[T]) Count() int64 {
	return count[T](f)
}

func (f *fluxImpl[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *fluxImpl[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *fluxImpl[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *fluxImpl[T]) Emit() opt.Opt[T] {
	d, ok := <-f.data
	if ok {
		return opt.New(d)
	}
	return opt.Empty[T]()
}

func (f *fluxImpl[T]) Iterator() iterator.Iterator[T] {
	return newFluxIterator(f.data)
}

func (f *fluxImpl[T]) Close() {
	close(f.data)
}

type filteredFlux[T any] struct {
	origin Flux[T]
	filter fn.Predicate[T]
}

func (f *filteredFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *filteredFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *filteredFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *filteredFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *filteredFlux[T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *filteredFlux[T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *filteredFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *filteredFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *filteredFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *filteredFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *filteredFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *filteredFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *filteredFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *filteredFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *filteredFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *filteredFlux[T]) Emit() opt.Opt[T] {
	for {
		inner := f.origin.Emit()
		if !inner.IsPresent() {
			return opt.Empty[T]()
		} else {
			t := inner.Get()
			if f.filter(t) {
				return inner
			}
		}
	}
}

func (f *filteredFlux[T]) Iterator() iterator.Iterator[T] {
	pipe := make(chan T)
	go func() {
		itr := f.origin.Iterator()
		for itr.HasNext() {
			n := itr.Next()
			if f.filter(n) {
				pipe <- n
			}
		}
		close(pipe)
	}()
	return newFluxIterator(pipe)
}

func (f *filteredFlux[T]) Close() {
	f.origin.Close()
}

type sortedFlux[T any] struct {
	origin     Flux[T]
	comparator fn.Comparator[T]
	iterator   iterator.Iterator[T]
}

func (f *sortedFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *sortedFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *sortedFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *sortedFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *sortedFlux[T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *sortedFlux[T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *sortedFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *sortedFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		next := itr.Next()
		consumer(next)
	}
}

func (f *sortedFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *sortedFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *sortedFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *sortedFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *sortedFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *sortedFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *sortedFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *sortedFlux[T]) Emit() opt.Opt[T] {
	itr := f.Iterator()
	if itr.HasNext() {
		return opt.New(itr.Next())
	}
	return opt.Empty[T]()
}

func (f *sortedFlux[T]) Iterator() iterator.Iterator[T] {
	if f.iterator != nil {
		return f.iterator
	}
	var sorted []T
	itr := f.origin.Iterator()
	for itr.HasNext() {
		if len(sorted) == 0 {
			sorted = append(sorted, itr.Next())
		} else {
			n := itr.Next()
			for i, t := range sorted {
				if f.comparator(n, t) >= 0 {
					if i == len(sorted)-1 {
						sorted = append(sorted, n)
					} else if f.comparator(n, sorted[i+1]) < 0 {
						sorted = append(sorted[:i+1], append([]T{n}, sorted[i+1:]...)...)
					}
				}
			}
		}
	}
	pipe := make(chan T)
	go func() {
		for _, t := range sorted {
			pipe <- t
		}
		close(pipe)
	}()
	return newFluxIterator(pipe)
}

func (f *sortedFlux[T]) Close() {
	f.origin.Close()
}

type distinctFlux[T any] struct {
	origin   Flux[T]
	comparer fn.Comparer[T]
	emitted  []T
}

func (f *distinctFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *distinctFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *distinctFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *distinctFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *distinctFlux[T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *distinctFlux[T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *distinctFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *distinctFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *distinctFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *distinctFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *distinctFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *distinctFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *distinctFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *distinctFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *distinctFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *distinctFlux[T]) Emit() opt.Opt[T] {
	itr := f.Iterator()
	if itr.HasNext() {
		return opt.New(itr.Next())
	}
	return opt.Empty[T]()
}

func (f *distinctFlux[T]) Iterator() iterator.Iterator[T] {
	pipe := make(chan T)
	go func() {
		itr := f.origin.Iterator()
		for itr.HasNext() {
			next := itr.Next()
			hit := false
			for _, t := range f.emitted {
				if f.comparer(t, next) {
					hit = true
					break
				}
			}
			if !hit {
				f.emitted = append(f.emitted, next)
				pipe <- next
			}
		}
		close(pipe)
	}()
	return newFluxIterator(pipe)
}

func (f *distinctFlux[T]) Close() {
	f.origin.Close()
}

type consumerAwareFlux[T any] struct {
	origin   Flux[T]
	consumer fn.Consumer[T]
}

func (f *consumerAwareFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *consumerAwareFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *consumerAwareFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *consumerAwareFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *consumerAwareFlux[T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *consumerAwareFlux[T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *consumerAwareFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *consumerAwareFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *consumerAwareFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *consumerAwareFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *consumerAwareFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *consumerAwareFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *consumerAwareFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *consumerAwareFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *consumerAwareFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *consumerAwareFlux[T]) Emit() opt.Opt[T] {
	next := f.origin.Emit()
	if next.IsPresent() {
		f.consumer(next.Get())
		return next
	}
	return opt.Empty[T]()
}

func (f *consumerAwareFlux[T]) Iterator() iterator.Iterator[T] {
	itr := f.Iterator()
	pipe := make(chan T)
	go func() {
		for itr.HasNext() {
			next := itr.Next()
			pipe <- next
			f.consumer(next)
		}
	}()
	return newFluxIterator(pipe)
}

func (f *consumerAwareFlux[T]) Close() {
	f.origin.Close()
}

type limitedFlux[T any] struct {
	origin  Flux[T]
	maxSize int64
	count   int64
}

func (f *limitedFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *limitedFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *limitedFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *limitedFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *limitedFlux[T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *limitedFlux[T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *limitedFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *limitedFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		next := itr.Next()
		consumer(next)
	}
}

func (f *limitedFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *limitedFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *limitedFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *limitedFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *limitedFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *limitedFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *limitedFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *limitedFlux[T]) Emit() opt.Opt[T] {
	if f.count >= f.maxSize {
		return opt.Empty[T]()
	}
	f.count++
	return f.origin.Emit()
}

func (f *limitedFlux[T]) Iterator() iterator.Iterator[T] {
	itr := f.origin.Iterator()
	pipe := make(chan T)
	go func() {
		for itr.HasNext() {
			next := itr.Next()
			if f.count < f.maxSize {
				f.count++
				pipe <- next
			}
		}
		close(pipe)
	}()
	return newFluxIterator(pipe)
}

func (f *limitedFlux[T]) Close() {
	f.origin.Close()
}

type skippedFlux[T any] struct {
	origin Flux[T]
	skip   int64
	count  int64
}

func (f *skippedFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *skippedFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *skippedFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *skippedFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *skippedFlux[T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *skippedFlux[T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *skippedFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *skippedFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *skippedFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *skippedFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *skippedFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *skippedFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *skippedFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *skippedFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *skippedFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *skippedFlux[T]) Emit() opt.Opt[T] {
	for f.count < f.skip {
		_ = f.origin.Emit()
		f.count++
	}
	return f.origin.Emit()
}

func (f *skippedFlux[T]) Iterator() iterator.Iterator[T] {
	itr := f.origin.Iterator()
	pipe := make(chan T)
	go func() {
		for itr.HasNext() {
			next := itr.Next()
			f.count++
			if f.count <= f.skip {
				continue
			}
			pipe <- next
		}
		close(pipe)
	}()
	return newFluxIterator(pipe)
}

func (f *skippedFlux[T]) Close() {
	f.origin.Close()
}

type closeHookFlux[T any] struct {
	origin    Flux[T]
	closeHook func()
}

func (f *closeHookFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *closeHookFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *closeHookFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *closeHookFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *closeHookFlux[T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *closeHookFlux[T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *closeHookFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *closeHookFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *closeHookFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *closeHookFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *closeHookFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *closeHookFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *closeHookFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *closeHookFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *closeHookFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *closeHookFlux[T]) Emit() opt.Opt[T] {
	return f.origin.Emit()
}

func (f *closeHookFlux[T]) Iterator() iterator.Iterator[T] {
	return f.origin.Iterator()
}

func (f *closeHookFlux[T]) Close() {
	f.closeHook()
	f.origin.Close()
}

type mappedFlux[R, T any] struct {
	origin Flux[R]
	mapper fn.Function[R, T]
}

func (f *mappedFlux[R, T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *mappedFlux[R, T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *mappedFlux[R, T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *mappedFlux[R, T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *mappedFlux[R, T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *mappedFlux[R, T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *mappedFlux[R, T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *mappedFlux[R, T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *mappedFlux[R, T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *mappedFlux[R, T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *mappedFlux[R, T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *mappedFlux[R, T]) Count() int64 {
	return count[T](f)
}

func (f *mappedFlux[R, T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *mappedFlux[R, T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *mappedFlux[R, T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *mappedFlux[R, T]) Emit() opt.Opt[T] {
	r := f.origin.Emit()
	if !r.IsPresent() {
		return opt.Empty[T]()
	}
	t := f.mapper(r.Get())
	return opt.New(t)
}

func (f *mappedFlux[R, T]) Iterator() iterator.Iterator[T] {
	itr := f.origin.Iterator()
	pipe := make(chan T)
	go func() {
		for itr.HasNext() {
			r := itr.Next()
			t := f.mapper(r)
			pipe <- t
		}
		close(pipe)
	}()
	return newFluxIterator(pipe)
}

func (f *mappedFlux[R, T]) Close() {
	f.origin.Close()
}

type flatMappedFlux[R, T any] struct {
	origin Flux[R]
	mapper fn.Function[R, Flux[T]]
}

func (f *flatMappedFlux[R, T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *flatMappedFlux[R, T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *flatMappedFlux[R, T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *flatMappedFlux[R, T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumerAwareFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *flatMappedFlux[R, T]) Limit(size int64) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *flatMappedFlux[R, T]) Skip(n int64) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *flatMappedFlux[R, T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *flatMappedFlux[R, T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *flatMappedFlux[R, T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *flatMappedFlux[R, T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *flatMappedFlux[R, T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *flatMappedFlux[R, T]) Count() int64 {
	return count[T](f)
}

func (f *flatMappedFlux[R, T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *flatMappedFlux[R, T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *flatMappedFlux[R, T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *flatMappedFlux[R, T]) Emit() opt.Opt[T] {
	r := f.origin.Emit()
	if !r.IsPresent() {
		return opt.Empty[T]()
	}
	return f.mapper(r.Get()).Emit()
}

func (f *flatMappedFlux[R, T]) Iterator() iterator.Iterator[T] {
	itr := f.origin.Iterator()
	pipe := make(chan T)
	go func() {
		for itr.HasNext() {
			next := itr.Next()
			nf := f.mapper(next)
			nitr := nf.Iterator()
			if nitr.HasNext() {
				nnext := nitr.Next()
				pipe <- nnext
			}
		}
		close(pipe)
	}()
	return newFluxIterator(pipe)
}

func (f *flatMappedFlux[R, T]) Close() {
	f.origin.Close()
}

func min[T any](f Flux[T], comparator fn.Comparator[T]) opt.Opt[T] {
	itr := f.Iterator()
	var min T
	if itr.HasNext() {
		min = itr.Next()
	} else {
		return opt.Empty[T]()
	}

	for itr.HasNext() {
		next := itr.Next()
		if comparator(next, min) < 0 {
			min = next
		}
	}
	return opt.New(min)
}

func max[T any](f Flux[T], comparator fn.Comparator[T]) opt.Opt[T] {
	itr := f.Iterator()
	var max T
	if itr.HasNext() {
		max = itr.Next()
	} else {
		return opt.Empty[T]()
	}

	for itr.HasNext() {
		next := itr.Next()
		if comparator(next, max) > 0 {
			max = next
		}
	}
	return opt.New(max)
}

func count[T any](f Flux[T]) int64 {
	var cnt int64
	f.ForEach(func(t T) {
		cnt++
	})
	return cnt
}

func anyMatch[T any](f Flux[T], predicate fn.Predicate[T]) bool {
	itr := f.Iterator()
	for itr.HasNext() {
		next := itr.Next()
		if predicate(next) {
			return true
		}
	}
	return false
}

func allMatch[T any](f Flux[T], predicate fn.Predicate[T]) bool {
	itr := f.Iterator()
	for itr.HasNext() {
		next := itr.Next()
		if !predicate(next) {
			return false
		}
	}
	return true
}

func noneMatch[T any](f Flux[T], predicate fn.Predicate[T]) bool {
	itr := f.Iterator()
	for itr.HasNext() {
		next := itr.Next()
		if predicate(next) {
			return false
		}
	}
	return true
}

func toSlice[T any](f Flux[T]) []T {
	var ret []T
	f.ForEach(func(t T) {
		ret = append(ret, t)
	})
	return ret
}
