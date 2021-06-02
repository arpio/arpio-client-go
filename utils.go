package arpio

// SliceContainsString checks if the string v is in the slice of strings s.
func SliceContainsString(v string, s []string) bool {
	for _, elem := range s {
		if v == elem {
			return true
		}
	}
	return false
}

// SliceContainsInt checks if the int v is in the slice of ints s.
func SliceContainsInt(v int, s []int) bool {
	for _, elem := range s {
		if v == elem {
			return true
		}
	}
	return false
}
