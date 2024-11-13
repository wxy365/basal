package fn

import "testing"

func TestCounter(t *testing.T) {
	var cnt Counter
	if cnt.Incr() != 1 {
		t.Fail()
	}
	if cnt.GetAndIncr() != 1 {
		t.Fail()
	}
	if cnt.Incr() != 3 {
		t.Fail()
	}
	cnt.Set(5)
	if cnt.GetAndIncr() != 5 {
		t.Fail()
	}
}
