package iterator

import "github.com/wxy365/basal/fn"

func Sort[T any](itr Iterator[T], comparator fn.Comparator[T]) Iterator[T] {
	s := make([]T, 0)
	for itr.HasNext() {
		next := itr.Next()
		for i, v := range s {
			if comparator(v, next) > 0 {
				t := s[i:]
				s = append(s[:i], next)
				s = append(s, t...)
				break
			}
		}

	}
	return OfSlice(s)
}
