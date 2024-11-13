package types

func Num2Bool[T BasicNumberUnion](t Number[T]) bool {
	return t.Bool()
}

func Num2Str[T BasicNumberUnion](t Number[T]) string {
	return t.String()
}

func Str2Num[T BasicNumberUnion](s string, n Number[T]) (T, error) {
	err := n.Parse(s)
	if err != nil {
		var t T
		return t, err
	}
	return n.UnBox(), nil
}
