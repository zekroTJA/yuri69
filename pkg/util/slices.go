package util

func IndexOf[T comparable](s []T, v T) int {
	if len(s) == 0 {
		return -1
	}

	for i, c := range s {
		if c == v {
			return i
		}
	}

	return -1
}

func Contains[T comparable](s []T, v T) bool {
	return IndexOf(s, v) != -1
}

func ContainsAny[T comparable](s []T, v []T) bool {
	if len(s) == 0 || len(v) == 0 {
		return false
	}

	for _, c := range v {
		if Contains(s, c) {
			return true
		}
	}

	return false
}

func ContainsAll[T comparable](s []T, v []T) bool {
	if len(s) == 0 && len(v) == 0 {
		return true
	}
	if len(s) == 0 {
		return false
	}

	for _, c := range v {
		if !Contains(s, c) {
			return false
		}
	}

	return true
}
