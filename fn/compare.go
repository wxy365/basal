package fn

type Comparator[T any] func(left, right T) int

func Compare[T any](left, right T, comparator Comparator[T]) int {
	return comparator(left, right)
}

func CompareWithKey[T, R any](left, right T, comparator Comparator[R], keyExtractor Function[T, R]) int {
	return comparator(keyExtractor(left), keyExtractor(right))
}

func ComparatorWithKey[T, R any](comparator Comparator[R], keyExtractor Function[T, R]) Comparator[T] {
	return func(left, right T) int {
		return CompareWithKey(left, right, comparator, keyExtractor)
	}
}

func BiCompareWithKey[T, U, R any](left T, right U, comparator Comparator[R], leftKeyExtractor Function[T, R], rightKeyExtractor Function[U, R]) int {
	return comparator(leftKeyExtractor(left), rightKeyExtractor(right))
}

// Comparer is a function used to determine whether two objects are equal
type Comparer[T any] func(left, right T) bool
