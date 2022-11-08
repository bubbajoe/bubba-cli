package util

func SliceMap[T any, G any](slice []T, mapFunc func(T) G) []G {
	newSlice := make([]G, len(slice))
	for s := range slice {
		newSlice[s] = mapFunc(slice[s])
	}
	return newSlice
}

func SliceFilter[T any, G any](slice []T, mapFunc func(T) *G) []G {
	newSlice := make([]G, len(slice))
	for s := range slice {
		g := mapFunc(slice[s])
		if g != nil {
			newSlice[s] = *g
		}
	}
	return newSlice
}

func MaptoSlice[K comparable, V any, G any](dict map[K]V, mapFunc func(K, V) G) []G {
	newSlice := make([]G, len(dict))
	index := 0
	for k, v := range dict {
		newSlice[index] = mapFunc(k, v)
		index++
	}
	return newSlice
}
func ChanToSlice[T any](chv <-chan T) []T {
	slv := make([]T, 0)
	for {
		v, ok := <-chv
		if !ok {
			return slv
		}
		slv = append(slv, v)
	}
}

func MergeSlices[T any](slices ...[]T) []T {
	newSlice := make([]T, 0)
	for _, slice := range slices {
		newSlice = append(newSlice, slice...)
	}
	return newSlice
}

func S(s string) *string {
	return &s
}

func ParseCommand(s string, delim byte) []string {
	var result []string
	inquote := false
	skip := false
	i := 0

	for j := 0; j < len(s); j++ {
		if skip {
			skip = false
			continue
		}
		c := s[j]
		if c == '\'' || c == '"' {
			inquote = !inquote
		} else if c == '\\' {
			skip = true
		} else if c == delim && !inquote {
			result = append(result, s[i:j])
			i = j + 1
		}
	}
	return append(result, s[i:])
}
