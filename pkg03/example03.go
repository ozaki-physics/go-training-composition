package pkg03

func Abs(n int) int {
	if n < 0 {
		return -1 * n
	} else {
		return n
	}
}

func ReverseAbs(n int) int {
	if n < 0 {
		return n
	} else {
		return -1 * n
	}
}
