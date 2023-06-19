package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type pairs []pair

type pair struct {
	amount int
	word   string
}

func (p pairs) toWordsSlice() []string {
	words := make([]string, len(p))

	for i, v := range p {
		words[i] = v.word
	}

	return words
}

func Top10(text string) []string {
	wordToAmountMap := map[string]int{}

	words := strings.Fields(text)
	for _, word := range words {
		wordToAmountMap[word]++
	}

	wordToAmountSlice := make(pairs, 0)

	for k, v := range wordToAmountMap {
		if v != 0 {
			wordToAmountSlice = append(wordToAmountSlice, pair{
				amount: v,
				word:   k,
			})
		}
	}

	sort.Slice(wordToAmountSlice, func(i, j int) bool {
		if wordToAmountSlice[i].amount == wordToAmountSlice[j].amount {
			return strings.Compare(wordToAmountSlice[i].word, wordToAmountSlice[j].word) != 1
		}

		return wordToAmountSlice[i].amount > wordToAmountSlice[j].amount
	})

	if len(wordToAmountSlice) > 10 {
		return wordToAmountSlice[:10].toWordsSlice()
	}

	return wordToAmountSlice.toWordsSlice()
}
