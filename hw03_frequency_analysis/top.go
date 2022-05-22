package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCounter = struct {
	Word  string
	Count int
}

func Top10(text string) []string {
	words := strings.Fields(text)

	if len(words) == 0 {
		return words
	}

	wordsWithCountMap := make(map[string]WordCounter)

	for _, word := range words {
		var currentCounter WordCounter
		if _, ok := wordsWithCountMap[word]; ok {
			currentCounter = wordsWithCountMap[word]
		} else {
			currentCounter = WordCounter{Word: word, Count: 0}
		}

		currentCounter.Count++

		wordsWithCountMap[word] = currentCounter
	}

	wordsWithCountSlice := make([]WordCounter, 0, len(wordsWithCountMap))
	for _, counter := range wordsWithCountMap {
		wordsWithCountSlice = append(wordsWithCountSlice, counter)
	}

	sort.Slice(wordsWithCountSlice, func(i, j int) bool {
		if wordsWithCountSlice[i].Count == wordsWithCountSlice[j].Count {
			return wordsWithCountSlice[i].Word < wordsWithCountSlice[j].Word
		}
		return wordsWithCountSlice[i].Count > wordsWithCountSlice[j].Count
	})
	wordsWithCountSlice = wordsWithCountSlice[0:10]

	result := make([]string, 0, 10)
	for _, counter := range wordsWithCountSlice {
		result = append(result, counter.Word)
	}

	return result
}
