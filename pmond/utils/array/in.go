package array

func In(arr []string, val string) bool {
	rel := map[string]int{}
	for _, s := range arr {
		rel[s] = 1
	}

	_, ok := rel[val]
	return ok
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
