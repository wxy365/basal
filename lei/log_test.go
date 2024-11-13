package lei

import (
	"testing"
)

func TestLog(t *testing.T) {
	e := New("{1}abc{0}", "eeee", "ffff")
	ErrorErrF("{0}fsdfasd", e, "fsd")
}
