package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Freq struct {
	Name    string
	Counter int
}

func Top10(s string) []string {
	// Place your code here.
	countMap := make(map[string]int)

	freqSlice := make([]Freq, 0)

	for _, s := range strings.Fields(s) {
		value := countMap[s]
		countMap[s] = value + 1
	}
	for key, value := range countMap {
		freqSlice = append(freqSlice, Freq{key, value})
	}

	sort.Slice(freqSlice, func(i, j int) bool {
		if freqSlice[i].Counter == freqSlice[j].Counter {
			return freqSlice[i].Name < freqSlice[j].Name
		}
		return freqSlice[i].Counter > freqSlice[j].Counter
	})

	resultSlice := make([]string, 0)
	for ind, value := range freqSlice {
		if ind == 10 {
			break
		}
		resultSlice = append(resultSlice, value.Name)
	}
	return resultSlice
}
