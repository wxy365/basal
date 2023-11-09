package opt

import (
	"fmt"
	"testing"
)

func TestOptional(t *testing.T) {
	var ex example
	o := New(ex)
	fmt.Println(o.Get())
}

type example struct {
	a string
	b int
}
