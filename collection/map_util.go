package coll

// MergeMap merges entries in map `from` into map `to`
func MergeMap[K comparable, V any](to, from map[K]V, replace ...bool) {
	r := false
	if len(replace) > 0 {
		r = replace[0]
	}
	if to == nil {
		panic("Merging map into nil")
	}
	if len(from) == 0 {
		return
	}
	if r {
		for k, v := range from {
			to[k] = v
		}
	} else {
		for k, v := range from {
			if _, exists := to[k]; !exists {
				to[k] = v
			}
		}
	}
}
