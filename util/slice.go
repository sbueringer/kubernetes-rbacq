package util

func Contains(stringArray []string, string string) bool {
	for _, a := range stringArray {
		if a == string {
			return true
		}
	}
	return false
}
