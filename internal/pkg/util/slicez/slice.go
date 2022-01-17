package slicez

// TODO(daniel-hutao): use set to improve the implementation

// SliceInSliceInterface is used to filter the items in s1 but not in s2
// Please use the SliceInSliceStr() or SliceInSliceInt() instead, they are more clear and easy to use.
func SliceInSliceInterface(s1, s2 interface{}) interface{} {
	slice1 := s1.([]interface{})
	slice2 := s2.([]interface{})
	retSlice := make([]interface{}, 0)
	for _, x := range slice1 {
		i := 0
		for ; i < len(slice2); i++ {
			if x == slice2[i] {
				break
			}
		}
		if i == len(slice2) {
			retSlice = append(retSlice, x)
		}
	}
	return retSlice
}

// SliceInSliceStr is used to filter the string items in s1 but not in s2
// eg:
// s1 -> ["a","b","c"]; s2 -> ["a","d","e"] => return ["b","c"]
// s1 -> ["a","b","c"]; s2 -> ["a","b","c"] => return []
func SliceInSliceStr(s1, s2 []string) []string {
	retSlice := make([]string, 0)
	for _, x := range s1 {
		i := 0
		for ; i < len(s2); i++ {
			if x == s2[i] {
				break
			}
		}
		if i == len(s2) {
			retSlice = append(retSlice, x)
		}
	}
	return retSlice
}

// SliceInSliceInt is used to filter the int items in s1 but not in s2
// eg:
// s1 -> [1,2,3]; s2 -> [1,5,6] => return [2,3]
// s1 -> [1,2,3]; s2 -> [1,2,3] => return []
func SliceInSliceInt(s1, s2 []int) []int {
	retSlice := make([]int, 0)
	for _, x := range s1 {
		i := 0
		for ; i < len(s2); i++ {
			if x == s2[i] {
				break
			}
		}
		if i == len(s2) {
			retSlice = append(retSlice, x)
		}
	}
	return retSlice
}
