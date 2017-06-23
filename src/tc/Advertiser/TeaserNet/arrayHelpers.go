package main

type StringArray []string

func (arr StringArray) Count() int {
	return len(arr)
}

func (arr StringArray) Contains(v string) bool {
	for _, b := range arr {
		if b == v {
			return true
		}
	}

	return false
}
