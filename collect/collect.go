package collect

func Contains(items []int, n int) bool {
	for _, v := range items {
		if v == n {
			return true
		}
	}

	return false
}

func Intersection(a, b []int) (c []int) {
	m := make(map[int]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func Union(a, b []int) []int {
	m := make(map[int]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			a = append(a, item)
		}
	}
	return a
}

func Difference(a, b []int) []int {
	m := make(map[int]bool)
	for _, item := range b {
		m[item] = true
	}

	res := []int{}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			res = append(res, item)
		}
	}

	return res
}
