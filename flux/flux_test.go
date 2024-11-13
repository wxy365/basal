package flux

import (
	"fmt"
	"testing"
)

var (
	s = []int{2, 4, 5, 23, 2, 54, 2, 54, 6, 95, 23, 54, 323, 94}
)

func TestFlux(t *testing.T) {
	//FromSlice(s).Filter(func(i int) bool {
	//	return i > 10
	//}).Sorted(func(left, right int) int {
	//	return left - right
	//}).Distinct(func(left, right int) bool {
	//	return left == right
	//}).Skip(2).Limit(1).ForEach(func(i int) {
	//	fmt.Println(i)
	//})

	FromSlice(s).Limit(2).ForEach(func(i int) {
		fmt.Println(i)
	})
}
