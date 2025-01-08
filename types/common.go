package types

func IsEmpty[T comparable](t T) bool {
	var zero T
	return zero == t
}
