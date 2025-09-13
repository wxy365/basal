package maps

import (
	"testing"
)

func TestMapToObj(t *testing.T) {
	m := map[string]any{
		"name":   "wxy",
		"age":    float64(18),
		"height": 1.85,
	}
	u := &User{}
	ToObj(m, &u)
	if u.Name1 != "wxy" {
		t.Fail()
	}
	if u.Age1 != 18 {
		t.Fail()
	}
	if u.Height1 != 1.85 {
		t.Fail()
	}
}

type User struct {
	Name1   string  `map:"name"`
	Age1    int     `map:"age"`
	Height1 float64 `map:"height"`
}
