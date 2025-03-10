package types

import (
	"bytes"
	"encoding/gob"
)

func IsEmpty[T comparable](t T) bool {
	var zero T
	return zero == t
}

func DeepClone[T any](t T) (T, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	if err := enc.Encode(t); err != nil {
		return t, err
	}
	var res T
	err := dec.Decode(&res)
	return res, err
}
