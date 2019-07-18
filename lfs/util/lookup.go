package util

type StringLookup []string
type IntLookup []int

type LookupTable []struct {
	returnVal int
	lookup    StringLookup
}

func (s LookupTable) Contains(x string) (bool, int) {
	for _, a := range s {
		for _, n := range a.lookup {
			if x == n {
				return true, a.returnVal
			}
		}
	}
	return false, 0
}

func (s StringLookup) Contains(x string) bool {
	for _, n := range s {
		if x == n {
			return true
		}
	}
	return false
}

func (s IntLookup) Contains(x int) bool {
	for _, n := range s {
		if x == n {
			return true
		}
	}
	return false
}
