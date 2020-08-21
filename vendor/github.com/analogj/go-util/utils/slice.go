package utils

func SliceIncludes(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

//func indexOf(answers []interface{}, item interface{}) (int) {
//	for k, v := range answers {
//		if v == item {
//			return k
//		}
//	}
//	return -1
//}
