package iterator

func Skip[T any](itr Iterator[T], skip uint) Iterator[T] {
	var count uint
	for itr.HasNext() && count < skip {
		itr.Next()
		count++
	}
	return itr
}
