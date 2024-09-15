package utils

func Max[K ~int](a K, b K) K {
	if a > b {
		return a
	}
	return b
}
