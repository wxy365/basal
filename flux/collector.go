package flux

import (
	"github.com/wxy365/basal/fn"
	"strings"
)

type Collector[T, A, R any] interface {
	Supplier() fn.Supplier[A]
	Gatherer() fn.BiConsumer[A, T]
	Synthesizer() fn.BiFunction[A, A, A]
	Finisher() fn.Function[A, R]
}

func NewCollector[T, A, R any](
	supplier fn.Supplier[A],
	gatherer fn.BiConsumer[A, T],
	synthesizer fn.BiFunction[A, A, A],
	finisher fn.Function[A, R]) Collector[T, A, R] {
	return &collectorImpl[T, A, R]{
		supplier:    supplier,
		gatherer:    gatherer,
		synthesizer: synthesizer,
		finisher:    finisher,
	}
}

func JoinStringCollector(delimiter, prefix, suffix string) Collector[string, *[]string, string] {
	return &collectorImpl[string, *[]string, string]{
		supplier: func() *[]string {
			ret := make([]string, 0)
			return &ret
		},
		gatherer: func(t *[]string, u string) {
			s := *t
			s = append(s, u)
		},
		synthesizer: nil,
		finisher: func(i *[]string) string {
			s := *i
			return prefix + strings.Join(s, delimiter) + suffix
		},
	}
}

func MergeMapCollector[K comparable, V any]() Collector[map[K]V, map[K]V, map[K]V] {
	return &collectorImpl[map[K]V, map[K]V, map[K]V]{
		supplier: func() map[K]V {
			return make(map[K]V)
		},
		gatherer: func(t map[K]V, u map[K]V) {
			for k, v := range u {
				t[k] = v
			}
		},
		synthesizer: nil,
		finisher: func(m map[K]V) map[K]V {
			return m
		},
	}
}

func SliceCollector[T any]() Collector[T, *[]T, []T] {
	return &collectorImpl[T, *[]T, []T]{
		supplier: func() *[]T {
			s := make([]T, 0)
			return &s
		},
		gatherer: func(t *[]T, u T) {
			s := *t
			s = append(s, u)
			*t = s
		},
		synthesizer: nil,
		finisher: func(i *[]T) []T {
			return *i
		},
	}
}

type collectorImpl[T, A, R any] struct {
	supplier    fn.Supplier[A]
	gatherer    fn.BiConsumer[A, T]
	synthesizer fn.BiFunction[A, A, A]
	finisher    fn.Function[A, R]
}

func (c *collectorImpl[T, A, R]) Supplier() fn.Supplier[A] {
	return c.supplier
}

func (c *collectorImpl[T, A, R]) Gatherer() fn.BiConsumer[A, T] {
	return c.gatherer
}

func (c *collectorImpl[T, A, R]) Synthesizer() fn.BiFunction[A, A, A] {
	return c.synthesizer
}

func (c *collectorImpl[T, A, R]) Finisher() fn.Function[A, R] {
	return c.finisher
}
