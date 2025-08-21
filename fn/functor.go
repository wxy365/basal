package fn

type Supplier[T any] func() T

type Predicate[T any] func(T) bool

func (p Predicate[T]) And(other Predicate[T]) Predicate[T] {
	if other == nil {
		return p
	}
	return func(t T) bool {
		return p(t) && other(t)
	}
}

func (p Predicate[T]) Negate() Predicate[T] {
	return func(t T) bool {
		return !p(t)
	}
}

func (p Predicate[T]) Or(other Predicate[T]) Predicate[T] {
	if other == nil {
		return p
	}
	return func(t T) bool {
		return p(t) || other(t)
	}
}

type BiPredicate[T, U any] func(t T, u U) bool

func (b BiPredicate[T, U]) And(other BiPredicate[T, U]) BiPredicate[T, U] {
	if other == nil {
		return b
	}
	return func(t T, u U) bool {
		return b(t, u) && other(t, u)
	}
}

func (b BiPredicate[T, U]) Negate() BiPredicate[T, U] {
	return func(t T, u U) bool {
		return !b(t, u)
	}
}

func (b BiPredicate[T, U]) Or(other BiPredicate[T, U]) BiPredicate[T, U] {
	if other == nil {
		return b
	}
	return func(t T, u U) bool {
		return b(t, u) || other(t, u)
	}
}

type Function[T, R any] func(T) R

type TryFunction[T, R any] func(T) (R, bool)

func ComposeFunction[V, T, R any](cur Function[T, R], before Function[V, T]) Function[V, R] {
	return func(v V) R {
		return cur(before(v))
	}
}

func NoOpsFunction[T any]() Function[T, T] {
	return func(t T) T {
		return t
	}
}

type BiFunction[T, U, R any] func(t T, u U) R

func ComposeBiFunction[T, U, R, V any](cur BiFunction[T, U, R], after Function[R, V]) BiFunction[T, U, V] {
	return func(t T, u U) V {
		return after(cur(t, u))
	}
}

type Consumer[T any] func(T)

func (c Consumer[T]) AndThen(after Consumer[T]) Consumer[T] {
	if after == nil {
		return c
	}
	return func(t T) {
		c(t)
		after(t)
	}
}

type BiConsumer[T, U any] func(t T, u U)

func (b BiConsumer[T, U]) AndThen(after BiConsumer[T, U]) BiConsumer[T, U] {
	if after == nil {
		return b
	}
	return func(t T, u U) {
		b(t, u)
		after(t, u)
	}
}
