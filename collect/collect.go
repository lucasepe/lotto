package collect

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lucasepe/lotto/data"
)

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

func Count(recs []data.Record, id string, debug bool) map[int]int {
	res := map[int]int{}
	for i := 1; i <= 90; i++ {
		res[i] = 0
	}

	for _, el := range recs {

		for key, vals := range el.Numbers {
			if strings.EqualFold(key, id) {
				if debug {
					fmt.Fprintf(os.Stderr, "%d - %s: %v\n", el.Day, key, vals)
				}
				for _, n := range vals {
					res[n] = res[n] + 1
				}
			}
		}
	}

	return res
}

func CountIf(recs []data.Record, id string, accept func(int, int) bool) (map[int]int, []int) {
	res := map[int]int{}

	counters := map[int]int{}
	for _, el := range recs {
		for key, vals := range el.Numbers {
			if strings.EqualFold(key, id) {
				for _, n := range vals {
					counters[n] = counters[n] + 1

					if accept(n, counters[n]) {
						res[n] = counters[n]
					}
				}
			}
		}
	}

	return res, KeysOnly(res, true)
}

func KeysOnly(src map[int]int, sorted bool) []int {
	res := make([]int, 0, len(src))

	for k := range src {
		res = append(res, k)
	}

	if sorted {
		sort.Ints(res)
	}

	return res
}
