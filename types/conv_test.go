package types

import (
	"testing"
)

func TestParseString(t *testing.T) {
	str := "1"
	num, err := Str2Num[int](str, new(Int))
	if err != nil {
		return
	}
	if num != 1 {
		t.Fail()
	}
}
