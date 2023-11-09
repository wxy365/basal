package coll

func LookupInSlice[T comparable](col []T, tgt T) int {
	for i, item := range col {
		if item == tgt {
			return i
		}
	}
	return -1
}

func SliceEqIgnoreOrder[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	idxes := NewIndexSlice(len(s2))
	for i1 := 0; i1 < len(s1); {
		item1 := s1[i1]
		found := false
		for j, i2 := range idxes {
			if item1 == s2[i2] {
				found = true
				idxes = append(idxes[:j], idxes[j+1:]...)
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func NewIndexSlice(n int) []int {
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = i
	}
	return ret
}
