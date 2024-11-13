package rflt

import (
	"reflect"
	"testing"
)

func TestSetFieldValue(t *testing.T) {

	newVal := "abc"

	f := foo{}

	err := SetFieldValue[string](&f, "Bar", newVal)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if f.Bar != newVal {
		t.Fail()
	}

	err = SetFieldValueAny(&f, "Bar", newVal)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if f.Bar != newVal {
		t.Fail()
	}

	err = SetFieldValue[*string](&f, "Bar", &newVal)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if f.Bar != newVal {
		t.Fail()
	}

	b := bar{}
	err = SetFieldValue[string](&b, "Foo", newVal)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if *b.Foo != newVal {
		t.Fail()
	}

	err = SetFieldValue[*string](&b, "Foo", &newVal)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	if *b.Foo != newVal {
		t.Fail()
	}
}

func TestPopulateValueFromString(t *testing.T) {
	var a string
	err := UnmarshalValue(reflect.ValueOf(&a), "aaa")
	if err != nil {
		panic(err)
	}
	if a != "aaa" {
		t.Fail()
	}
}

type foo struct {
	Bar string
}
type bar struct {
	Foo *string
}
