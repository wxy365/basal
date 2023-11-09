package lei

import "testing"

func TestLog(t *testing.T) {
	e := New("{1}abc{0}", "eeee", "ffff")
	//Panic("fasd;fka;sd")
	//ErrorErrF("{0}fsdfasd", nil, "fsd")
	DebugErr(e)
}
