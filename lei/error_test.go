package lei

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	err := New("{0} bar", "foo")
	err1 := Wrap("{0} {1}", errors.New("something bad happened"), "oh", "no")
	err2 := Wrap("{0} {1}", err, "come", "on")
	if err.Error() != "{\"msg\":\"foo bar\"}" {
		t.Fail()
	}
	if err1.Error() != "{\"msg\":\"oh no\",\"cause\":\"something bad happened\"}" {
		t.Fail()
	}
	if err2.Error() != "{\"msg\":\"come on\",\"cause\":{\"msg\":\"foo bar\"}}" {
		t.Fail()
	}
}
