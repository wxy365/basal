package iterator

import "github.com/wxy365/basal/fn"

type Iterator[T any] interface {
	HasNext() bool
	Next() T
	Remove() // 从底层容器中删除上一个返回的元素
	ForEach(consumer fn.Consumer[T])
}
