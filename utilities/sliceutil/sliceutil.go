package sliceutil

func ShiftEnd[T any](s []T, x int) []T {
	if len(s) <= 1 || x >= len(s) {
		return s
	}
	tmp := s[x]

	return append(append(s[:x], s[x+1:]...), tmp)
}
