package gocc

func runeCmp(r1 []rune, r2 []rune) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i, r := range r1 {
		if r2[i] != r {
			return false
		}
	}
	return true
}
