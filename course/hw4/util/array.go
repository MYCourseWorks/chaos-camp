package util

// IsInArray if a val is in array returns true
func IsInArray(val string, array []string) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}

	return false
}
