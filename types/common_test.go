package types

import (
	"testing"
)

func IsEmptyWithZeroValue(t *testing.T) {
	if !IsEmpty(0) {
		t.Fail()
	}
}

func IsEmptyWithNonZeroValue(t *testing.T) {
	if IsEmpty(1) {
		t.Fail()
	}
}

func TestDeepCloneWithValidInput(t *testing.T) {
	type TestStruct struct {
		Field1 int
		Field2 string
	}
	original := &TestStruct{Field1: 1, Field2: "test"}
	clone, err := DeepClone(original)
	if err != nil {
		t.Fail()
	}
	if clone != original {
		t.Fail()
	}
}

func TestDeepCloneWithInvalidInput(t *testing.T) {
	_, err := DeepClone(make(chan int))
	if err == nil {
		t.Fail()
	}
}
