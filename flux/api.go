package flux

import (
	"github.com/wxy365/basal/iterator"
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

func FromRange[T int](start, end T) Flux[T] {
	pipe := make(chan T)
	go func() {
		for i := start; i < end; i++ {
			pipe <- i
		}
		close(pipe)
	}()
	return &chanFlux[T]{pipe}
}
