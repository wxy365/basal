package slices

import "github.com/wxy365/basal/fn"

func NewIndexSlice(n int) []int {
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = i
	}
	return ret
}

func Lookup[T any](s []T, tgt T, comparer fn.Comparer[T]) int {
	for i, item := range s {
		if comparer(item, tgt) {
			return i
		}
	}
	return -1
}

func Contains[T any](s []T, t T, comparer fn.Comparer[T]) bool {
	return Lookup(s, t, comparer) >= 0
}

func Eq[T any](s []T, another []T, comparer fn.Comparer[T]) bool {
	if len(s) != len(another) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !comparer(s[i], another[i]) {
			return false
		}
	}
	return true
}

// EqIgOrder determines whether two slices are equal ignoring order
func EqIgOrder[T any](s []T, another []T, comparer fn.Comparer[T]) bool {
	if len(s) != len(another) {
		return false
	}
	idxes := NewIndexSlice(len(another))
	for i1 := 0; i1 < len(s); {
		item1 := s[i1]
		found := false
		for j, i2 := range idxes {
			if comparer(item1, another[i2]) {
				found = true
				idxes = append(idxes[:j], idxes[j+1:]...)
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func Merge[T any](s1 []T, s2 []T) []T {
	if len(s1) > len(s2) {
		for i, item := range s2 {
			s1[i] = item
		}
		return s1
	}
	for i, item := range s1 {
		s2[i] = item
	}
	return s2
}

func Del[T any](s []T, i int) []T {
	if i < 0 || i > len(s)-1 {
		return s
	}
	if i == 0 {
		return s[1:]
	}
	if i == len(s)-1 {
		return s[:i]
	}
	return append(s[:i], s[i+1:]...)
}

func DelRange[T any](s []T, from, to int) []T {
	if from < 0 {
		from = 0
	}
	if to > len(s) {
		to = len(s)
	}
	if from > to {
		return s
	}
	return append(s[:from], s[to:]...)
}

func Insert[T any](s []T, t T, i int) []T {
	if i < 0 {
		i = 0
	}
	if len(s) <= i {
		return append(s, t)
	}
	s = append(s[:i+1], s[i:]...)
	s[i] = t
	return s
}

func ForEach[T any](s []T, consumer fn.Consumer[T]) {
	for _, item := range s {
		consumer(item)
	}
}

func Map[T, R any](s []T, function fn.Function[T, R]) []R {
	res := make([]R, len(s))
	for i, item := range s {
		res[i] = function(item)
	}
	return res
}

func New[T any](t T, l int) []T {
	s := make([]T, l)
	s = Map(s, func(T) T {
		return t
	})
	return s
}
