package iterator

import (
	"fmt"
	"strings"
	"testing"
)

func TestMapIterator(t *testing.T) {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	m["d"] = 4

	var sb strings.Builder
	itr := OfMap(m)
	for itr.HasNext() {
		nxt := itr.Next()
		sb.WriteString(fmt.Sprintf("%s%d", nxt.GetKey(), nxt.GetValue()))
	}
	if sb.String() != "a1b2c3d4" {
		t.Fail()
	}
}
