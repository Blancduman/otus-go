package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	wordToAmountMap := map[string]int{}
	words := make([]string, 0)
	result := make([]string, 0)

	for _, v := range strings.Fields(strings.ToLower(text)) {
		word := strings.Trim(v, "\"!,.-")
		if word != "" {
			words = append(words, word)
		}
	}

	for _, word := range words {
		wordToAmountMap[word]++
	}

	sort.Slice(words, func(i, j int) bool {
		a := words[i]
		b := words[j]

		if wordToAmountMap[a] == wordToAmountMap[b] {
			return a < b
		}

		return wordToAmountMap[a] > wordToAmountMap[b]
	})

	for _, w := range words {
		if len(result) == 10 {
			break
		}

		canAdd := true
		for _, v := range result {
			if w == v {
				canAdd = false
				continue
			}
		}

		if canAdd {
			result = append(result, w)
		}
	}

	return result
}
