package ioc

import (
	"testing"
)

func TestReflect(t *testing.T) {
	a := &A{
		Field1: "a f1",
	}
	b := &B{
		Field1: "b f1",
	}
	c := &C{
		Field1: "c f1",
		Field2: "c f2",
	}
	d := &D{Field1: "d f1"}
	d2 := &D{Field1: "d2 f1"}

	Register(a)
	Register(b)
	Register(c)
	Register(d)
	Register(d2, "d2")

	if b.Field2 != a {
		t.Fail()
	}
	if b.Field3 != c {
		t.Fail()
	}
	if c.Field3 != d2 {
		t.Fail()
	}
}

type A struct {
	Field1 string
}

type B struct {
	Field1 string
	Field2 *A `autowired:""`
	Field3 *C `autowired:""`
}

type C struct {
	Field1 string
	Field2 string
	Field3 Iface `autowired:"d2"`
}

type D struct {
	Field1 string
}

func (d *D) Do() string {
	return d.Field1
}

type Iface interface {
	Do() string
}
