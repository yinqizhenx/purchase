package trick

import (
	"sort"
)

func cut(a []int, i, j int) []int {
	return append(a[:i], a[j+1:]...)
}

func reserve(a []int) []int {
	for l, r := 0, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return a
}

// moveToFront moves needle to the front of haystack, in place if possible.
func moveToFront(needle string, haystack []string) []string {
	if len(haystack) != 0 && haystack[0] == needle {
		return haystack
	}
	prev := needle
	for i, elem := range haystack {
		switch {
		case i == 0:
			haystack[0] = needle
			prev = elem
		case elem == needle:
			haystack[i] = prev
			return haystack
		default:
			haystack[i] = prev
			prev = elem
		}
	}
	return append(haystack, prev)
}

// deduplicate去重 In-place deduplicate (comparable)
func deduplicate(in []int) []int {
	// in := []int{3,2,1,4,3,2,1,4,1} // any item can be sorted
	sort.Ints(in)
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		in[j] = in[i]
	}
	result := in[:j+1]
	return result
}

func slidingWindow(size int, input []int) [][]int {
	// returns the input slice as the first element
	if len(input) <= size {
		return [][]int{input}
	}

	// allocate slice at the precise size we need
	r := make([][]int, 0, len(input)-size+1)

	for i, j := 0, size; j <= len(input); i, j = i+1, j+1 {
		r = append(r, input[i:j])
	}

	return r
}

// filter Filtering without allocating
func filter[T any](fn func(e T) bool, arr []T) []T {
	rst := arr[:0]
	for _, e := range arr {
		if fn(e) {
			rst = append(rst, e)
		}
	}
	return rst
}

func Insert[T any](arr []T, i int, val T) []T {
	arr = append(arr, val)
	copy(arr[i+1:], arr[i:])
	arr[i] = val
	return arr
}
