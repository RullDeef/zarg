package utils

func FindFirstOrPanic(len int, f func(int) bool) int {
	for i := 0; i < len; i++ {
		if f(i) {
			return i
		}
	}
	panic("element not found")
}
