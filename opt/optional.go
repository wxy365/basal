package opt

import (
	"github.com/wxy365/basal/fn"
	"reflect"
)

type Opt[T any] interface {
	Get() T
	IsPresent() bool
	IfPresent(c fn.Consumer[T])
	Filter(p fn.Predicate[T]) Opt[T]
	OrElse(other T) T
	OrElseGet(other fn.Supplier[T]) T
	OrElseErr(errSupplier fn.Supplier[error]) (T, error)
}

func Map[T, R any](o Opt[T], f fn.Function[T, R]) Opt[R] {
	if !o.IsPresent() {
		return Empty[R]()
	}
	return Of[R](f(o.Get()))
}

func FlatMap[T, R any](o Opt[T], f fn.Function[T, Opt[R]]) Opt[R] {
	if !o.IsPresent() {
		return Empty[R]()
	}
	return f(o.Get())
}

type optImpl[T any] struct {
	value *T
}

func (o *optImpl[T]) Get() T {
	return *o.value
}

func (o *optImpl[T]) IsPresent() bool {
	return o.value != nil
}

func (o *optImpl[T]) IfPresent(c fn.Consumer[T]) {
	if o.IsPresent() {
		c(*o.value)
	}
}

func (o *optImpl[T]) Filter(p fn.Predicate[T]) Opt[T] {
	if !o.IsPresent() {
		return o
	} else if p(*o.value) {
		return o
	} else {
		return Empty[T]()
	}
}

func (o *optImpl[T]) OrElse(other T) T {
	if o.IsPresent() {
		return *o.value
	}
	return other
}

func (o *optImpl[T]) OrElseGet(other fn.Supplier[T]) T {
	if o.IsPresent() {
		return *o.value
	}
	return other()
}

func (o *optImpl[T]) OrElseErr(errSupplier fn.Supplier[error]) (T, error) {
	if o.IsPresent() {
		return *o.value, nil
	}
	var r T
	return r, errSupplier()
}

func Empty[T any]() Opt[T] {
	return &optImpl[T]{}
}

func Of[T any](value T) Opt[T] {
	if reflect.ValueOf(value).IsZero() {
		return Empty[T]()
	}
	return &optImpl[T]{value: &value}
}
