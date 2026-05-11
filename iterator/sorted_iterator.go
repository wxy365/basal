package iterator

import (
	"sort"

	"github.com/wxy365/basal/fn"
)

func Sort[T any](itr Iterator[T], comparator fn.Comparator[T]) Iterator[T] {
	s := make([]T, 0)
	for itr.HasNext() {
		s = append(s, itr.Next())
	}
	sort.Slice(s, func(i, j int) bool {
		return comparator(s[i], s[j]) < 0
	})
	return FromSlice(s)
}
