package env

import (
	"os"
	"testing"
)

func TestGetObj(t *testing.T) {
	os.Setenv("INT_ARRAY", "[1,2,3,4]")
	r, e := GetObj[[]int64]("INT_ARRAY")
	if e != nil {
		t.Error(e)
	}
	expect := []int64{1, 2, 3, 4}
	if len(r) != len(expect) ||
		r[0] != expect[0] ||
		r[1] != expect[1] ||
		r[2] != expect[2] ||
		r[3] != expect[3] {
		t.Fail()
	}
}
