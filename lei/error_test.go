package lei

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := New("{0} is down", "a")
	err1 := Wrap("{0} is larger than {1}", fmt.Errorf("wuwuw"), "lll", "mmmm")
	err2 := Wrap("{0} is less than {1}", err, "kkk", "jjjj")
	fmt.Println(err.Error())
	fmt.Println(err1.Error())
	fmt.Println(err2.Error())
}
