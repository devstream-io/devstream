package slicez

// SliceInSlice is used to filter the items in s1 but not in s2
// eg1:
// s1 -> [1,2,3]; s2 -> [1,5,6] => return [2,3]
// s1 -> [1,2,3]; s2 -> [1,2,3] => return []
// eg2:
// s1 -> ["a","b","c"]; s2 -> ["a","d","e"] => return ["b","c"]
// s1 -> ["a","b","c"]; s2 -> ["a","b","c"] => return []
func SliceInSlice(s1, s2 interface{}) interface{} {
	slice1 := s1.([]interface{})
	slice2 := s2.([]interface{})
	retSlice := make([]interface{}, 0)
	for _, x := range slice1 {
		i := 0
		for ; i < len(slice2); i++ {
			if x == slice2[i] {
				continue
			}
		}
		if i == len(slice2) {
			retSlice = append(retSlice, x)
		}
	}
	return retSlice
}
