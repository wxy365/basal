package maps

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/wxy365/basal/text"
)

// Merge merges entries in map `another` into map `m`
func Merge[K comparable, V any](m map[K]V, another map[K]V, replace ...bool) map[K]V {
	r := false
	if len(replace) > 0 {
		r = replace[0]
	}
	if len(another) == 0 {
		return m
	}
	if r {
		for k, v := range another {
			m[k] = v
		}
	} else {
		for k, v := range another {
			if _, exists := m[k]; !exists {
				m[k] = v
			}
		}
	}
	return m
}

type MapOrSlice interface {
	map[string]any | []any
}

func ToObj[T any, M MapOrSlice](m M, t *T) {
	if len(m) == 0 {
		return
	}
	dcfg := mapstructure.DecoderConfig{
		ErrorUnset:           false,
		ErrorUnused:          false,
		ZeroFields:           false,
		WeaklyTypedInput:     true,
		Squash:               true,
		Result:               t,
		TagName:              "map",
		IgnoreUntaggedFields: false,
		MatchName: func(mapKey, fieldName string) bool {
			if strings.EqualFold(mapKey, fieldName) {
				return true
			}
			if strings.EqualFold(text.Snake2Pascal(mapKey), text.Snake2Pascal(fieldName)) {
				return true
			}
			return false
		},
	}
	decoder, _ := mapstructure.NewDecoder(&dcfg)
	_ = decoder.Decode(m)
}
