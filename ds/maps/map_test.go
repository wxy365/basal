package maps

import (
	"fmt"
	"testing"
)

func TestMapToObj(t *testing.T) {
	m := map[string]any{
		"Name":   "wxy",
		"Age":    float64(18),
		"Height": 1.85,
	}
	u := &User{}
	ToObj[*User](m, &u)
	fmt.Println(u)
}

type User struct {
	Name   string `json:"name"`
	Age    int
	Height float64
}
