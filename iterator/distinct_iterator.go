package iterator

import "github.com/wxy365/basal/fn"

func Distinct[T any](itr Iterator[T], comparer fn.Comparer[T]) Iterator[T] {
	s := make([]T, 0)
	for itr.HasNext() {
		next := itr.Next()
		exists := false
		for _, v := range s {
			if comparer(next, v) {
				exists = true
				break
			}
		}
		if !exists {
			s = append(s, next)
		}
	}
	return OfSlice(s)
}
