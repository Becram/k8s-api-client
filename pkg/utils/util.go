package util

func FindString(array []string, s string) bool {

	for _, v := range array {
		if v == s {
			return true
		}
	}
	return false

}
