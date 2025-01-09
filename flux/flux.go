package flux

import (
	"github.com/wxy365/basal/fn"
	"github.com/wxy365/basal/iterator"
	"github.com/wxy365/basal/opt"
	"github.com/wxy365/basal/types"
)

func FromMap[K comparable, V any](m map[K]V) Flux[iterator.MapEntry[K, V]] {
	pipe := make(chan iterator.MapEntry[K, V], len(m))
	go func() {
		for k, v := range m {
			entry := iterator.NewMapEntry(k, v)
			pipe <- entry
		}
		close(pipe)
	}()
	return &chanFlux[iterator.MapEntry[K, V]]{
		data: pipe,
	}
}

func FromSlice[T any](s []T) Flux[T] {
	return &sliceFlux[T]{tgt: s}
}

func FromRange[T types.BasicNumberUnion](start, end T) Flux[T] {
	pipe := make(chan T)
	go func() {
		for i := start; i < end; i++ {
			pipe <- i
		}
		close(pipe)
	}()
	return &chanFlux[T]{pipe}
}

type Flux[T any] interface {
	Filter(predicate fn.Predicate[T]) Flux[T]
	Sorted(comparator fn.Comparator[T]) Flux[T]
	Distinct(comparer fn.Comparer[T]) Flux[T]
	OnEmit(consumer fn.Consumer[T]) Flux[T]
	Limit(size uint) Flux[T]
	Skip(n uint) Flux[T]
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

func Map[T, R any](f Flux[T], mapper fn.Function[T, R]) Flux[R] {
	return &mappedFlux[T, R]{
		origin: f,
		mapper: mapper,
	}
}

func FlatMap[T, R any](f Flux[T], mapper fn.Function[T, Flux[R]]) Flux[R] {
	return &flatMappedFlux[T, R]{
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
	return &consumedFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *filteredFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *filteredFlux[T]) Skip(n uint) Flux[T] {
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
	return iterator.Filter(f.origin.Iterator(), f.filter)
}

func (f *filteredFlux[T]) Close() {
	f.origin.Close()
}

type sortedFlux[T any] struct {
	origin     Flux[T]
	comparator fn.Comparator[T]
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
	return &consumedFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *sortedFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *sortedFlux[T]) Skip(n uint) Flux[T] {
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
		return opt.Of(itr.Next())
	}
	return opt.Empty[T]()
}

func (f *sortedFlux[T]) Iterator() iterator.Iterator[T] {
	return iterator.Sort(f.origin.Iterator(), f.comparator)
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
	return &consumedFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *distinctFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *distinctFlux[T]) Skip(n uint) Flux[T] {
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
		return opt.Of(itr.Next())
	}
	return opt.Empty[T]()
}

func (f *distinctFlux[T]) Iterator() iterator.Iterator[T] {
	return iterator.Distinct(f.origin.Iterator(), f.comparer)
}

func (f *distinctFlux[T]) Close() {
	f.origin.Close()
}

type consumedFlux[T any] struct {
	origin   Flux[T]
	consumer fn.Consumer[T]
}

func (f *consumedFlux[T]) Filter(predicate fn.Predicate[T]) Flux[T] {
	return &filteredFlux[T]{
		origin: f,
		filter: predicate,
	}
}

func (f *consumedFlux[T]) Sorted(comparator fn.Comparator[T]) Flux[T] {
	return &sortedFlux[T]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *consumedFlux[T]) Distinct(comparer fn.Comparer[T]) Flux[T] {
	return &distinctFlux[T]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *consumedFlux[T]) OnEmit(consumer fn.Consumer[T]) Flux[T] {
	return &consumedFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *consumedFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *consumedFlux[T]) Skip(n uint) Flux[T] {
	return &skippedFlux[T]{
		origin: f,
		skip:   n,
	}
}

func (f *consumedFlux[T]) OnClose(hook func()) Flux[T] {
	return &closeHookFlux[T]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *consumedFlux[T]) ForEach(consumer fn.Consumer[T]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *consumedFlux[T]) ToSlice() []T {
	return toSlice[T](f)
}

func (f *consumedFlux[T]) Min(comparator fn.Comparator[T]) opt.Opt[T] {
	return min[T](f, comparator)
}

func (f *consumedFlux[T]) Max(comparator fn.Comparator[T]) opt.Opt[T] {
	return max[T](f, comparator)
}

func (f *consumedFlux[T]) Count() int64 {
	return count[T](f)
}

func (f *consumedFlux[T]) AnyMatch(predicate fn.Predicate[T]) bool {
	return anyMatch[T](f, predicate)
}

func (f *consumedFlux[T]) AllMatch(predicate fn.Predicate[T]) bool {
	return allMatch[T](f, predicate)
}

func (f *consumedFlux[T]) NoneMatch(predicate fn.Predicate[T]) bool {
	return noneMatch[T](f, predicate)
}

func (f *consumedFlux[T]) Emit() opt.Opt[T] {
	next := f.origin.Emit()
	if next.IsPresent() {
		f.consumer(next.Get())
		return next
	}
	return opt.Empty[T]()
}

func (f *consumedFlux[T]) Iterator() iterator.Iterator[T] {
	return iterator.Consume(f.origin.Iterator(), f.consumer)
}

func (f *consumedFlux[T]) Close() {
	f.origin.Close()
}

type limitedFlux[T any] struct {
	origin  Flux[T]
	maxSize uint
	count   uint
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
	return &consumedFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *limitedFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *limitedFlux[T]) Skip(n uint) Flux[T] {
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
	return iterator.Limit(itr, f.maxSize)
}

func (f *limitedFlux[T]) Close() {
	f.origin.Close()
}

type skippedFlux[T any] struct {
	origin Flux[T]
	skip   uint
	count  uint
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
	return &consumedFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *skippedFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *skippedFlux[T]) Skip(n uint) Flux[T] {
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
	return iterator.Skip[T](f.origin.Iterator(), f.skip)
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
	return &consumedFlux[T]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *closeHookFlux[T]) Limit(size uint) Flux[T] {
	return &limitedFlux[T]{
		origin:  f,
		maxSize: size,
	}
}

func (f *closeHookFlux[T]) Skip(n uint) Flux[T] {
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

type mappedFlux[T, R any] struct {
	origin Flux[T]
	mapper fn.Function[T, R]
}

func (f *mappedFlux[T, R]) Filter(predicate fn.Predicate[R]) Flux[R] {
	return &filteredFlux[R]{
		origin: f,
		filter: predicate,
	}
}

func (f *mappedFlux[T, R]) Sorted(comparator fn.Comparator[R]) Flux[R] {
	return &sortedFlux[R]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *mappedFlux[T, R]) Distinct(comparer fn.Comparer[R]) Flux[R] {
	return &distinctFlux[R]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *mappedFlux[T, R]) OnEmit(consumer fn.Consumer[R]) Flux[R] {
	return &consumedFlux[R]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *mappedFlux[T, R]) Limit(size uint) Flux[R] {
	return &limitedFlux[R]{
		origin:  f,
		maxSize: size,
	}
}

func (f *mappedFlux[T, R]) Skip(n uint) Flux[R] {
	return &skippedFlux[R]{
		origin: f,
		skip:   n,
	}
}

func (f *mappedFlux[T, R]) OnClose(hook func()) Flux[R] {
	return &closeHookFlux[R]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *mappedFlux[T, R]) ForEach(consumer fn.Consumer[R]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *mappedFlux[T, R]) ToSlice() []R {
	return toSlice[R](f)
}

func (f *mappedFlux[T, R]) Min(comparator fn.Comparator[R]) opt.Opt[R] {
	return min[R](f, comparator)
}

func (f *mappedFlux[T, R]) Max(comparator fn.Comparator[R]) opt.Opt[R] {
	return max[R](f, comparator)
}

func (f *mappedFlux[T, R]) Count() int64 {
	return count[R](f)
}

func (f *mappedFlux[T, R]) AnyMatch(predicate fn.Predicate[R]) bool {
	return anyMatch[R](f, predicate)
}

func (f *mappedFlux[T, R]) AllMatch(predicate fn.Predicate[R]) bool {
	return allMatch[R](f, predicate)
}

func (f *mappedFlux[T, R]) NoneMatch(predicate fn.Predicate[R]) bool {
	return noneMatch[R](f, predicate)
}

func (f *mappedFlux[T, R]) Emit() opt.Opt[R] {
	r := f.origin.Emit()
	if !r.IsPresent() {
		return opt.Empty[R]()
	}
	t := f.mapper(r.Get())
	return opt.Of(t)
}

func (f *mappedFlux[T, R]) Iterator() iterator.Iterator[R] {
	itr := f.origin.Iterator()
	return iterator.Map(itr, f.mapper)
}

func (f *mappedFlux[T, R]) Close() {
	f.origin.Close()
}

type flatMappedFlux[T, R any] struct {
	origin    Flux[T]
	mapper    fn.Function[T, Flux[R]]
	mappedEle Flux[R]
}

func (f *flatMappedFlux[T, R]) Filter(predicate fn.Predicate[R]) Flux[R] {
	return &filteredFlux[R]{
		origin: f,
		filter: predicate,
	}
}

func (f *flatMappedFlux[T, R]) Sorted(comparator fn.Comparator[R]) Flux[R] {
	return &sortedFlux[R]{
		origin:     f,
		comparator: comparator,
	}
}

func (f *flatMappedFlux[T, R]) Distinct(comparer fn.Comparer[R]) Flux[R] {
	return &distinctFlux[R]{
		origin:   f,
		comparer: comparer,
	}
}

func (f *flatMappedFlux[T, R]) OnEmit(consumer fn.Consumer[R]) Flux[R] {
	return &consumedFlux[R]{
		origin:   f,
		consumer: consumer,
	}
}

func (f *flatMappedFlux[T, R]) Limit(size uint) Flux[R] {
	return &limitedFlux[R]{
		origin:  f,
		maxSize: size,
	}
}

func (f *flatMappedFlux[T, R]) Skip(n uint) Flux[R] {
	return &skippedFlux[R]{
		origin: f,
		skip:   n,
	}
}

func (f *flatMappedFlux[T, R]) OnClose(hook func()) Flux[R] {
	return &closeHookFlux[R]{
		origin:    f,
		closeHook: hook,
	}
}

func (f *flatMappedFlux[T, R]) ForEach(consumer fn.Consumer[R]) {
	itr := f.Iterator()
	for itr.HasNext() {
		consumer(itr.Next())
	}
}

func (f *flatMappedFlux[T, R]) ToSlice() []R {
	return toSlice[R](f)
}

func (f *flatMappedFlux[T, R]) Min(comparator fn.Comparator[R]) opt.Opt[R] {
	return min[R](f, comparator)
}

func (f *flatMappedFlux[T, R]) Max(comparator fn.Comparator[R]) opt.Opt[R] {
	return max[R](f, comparator)
}

func (f *flatMappedFlux[T, R]) Count() int64 {
	return count[R](f)
}

func (f *flatMappedFlux[T, R]) AnyMatch(predicate fn.Predicate[R]) bool {
	return anyMatch[R](f, predicate)
}

func (f *flatMappedFlux[T, R]) AllMatch(predicate fn.Predicate[R]) bool {
	return allMatch[R](f, predicate)
}

func (f *flatMappedFlux[T, R]) NoneMatch(predicate fn.Predicate[R]) bool {
	return noneMatch[R](f, predicate)
}

func (f *flatMappedFlux[T, R]) Emit() opt.Opt[R] {
	if f.mappedEle == nil {
		r := f.origin.Emit()
		if r.IsPresent() {
			f.mappedEle = f.mapper(r.Get())
		} else {
			return opt.Empty[R]()
		}
	} else {
		t := f.mappedEle.Emit()
		if t.IsPresent() {
			return t
		} else {
			f.mappedEle = nil
		}
	}
	return f.Emit()
}

func (f *flatMappedFlux[T, R]) Iterator() iterator.Iterator[R] {
	itrMapper := func(r T) iterator.Iterator[R] {
		return f.mapper(r).Iterator()
	}
	return iterator.FlatMap[T, R](f.origin.Iterator(), itrMapper)
}

func (f *flatMappedFlux[T, R]) Close() {
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
	return opt.Of(min)
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
	return opt.Of(max)
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
