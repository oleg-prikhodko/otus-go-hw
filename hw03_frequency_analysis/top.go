package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	counter := make(map[string]int)
	for _, word := range strings.Fields(text) {
		counter[word]++
	}

	top := make([]string, 0, len(counter))
	for k := range counter {
		top = append(top, k)
	}

	sort.Slice(top, func(i, j int) bool {
		diff := counter[top[i]] - counter[top[j]]
		if diff == 0 {
			return top[i] < top[j]
		}

		return diff > 0
	})

	if len(top) > 10 {
		return top[:10]
	}

	return top
}
