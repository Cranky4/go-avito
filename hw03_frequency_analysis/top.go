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
	const maxResults = 10

	words := strings.Fields(text)

	if len(words) == 0 {
		return words
	}

	// дробим строку на подстроки
	wordsWithCountMap := make(map[string]WordCounter)

	// накручиваем мапу структур для подсчета, используя ключ мапы, как признак уже найденного слова
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

	// превращаем мапу в слайс, чтобы сортировать
	wordsWithCountSlice := make([]WordCounter, 0, len(wordsWithCountMap))
	for _, counter := range wordsWithCountMap {
		wordsWithCountSlice = append(wordsWithCountSlice, counter)
	}

	// сортируем по убыванию количества
	sort.Slice(wordsWithCountSlice, func(i, j int) bool {
		if wordsWithCountSlice[i].Count == wordsWithCountSlice[j].Count {
			return wordsWithCountSlice[i].Word < wordsWithCountSlice[j].Word
		}
		return wordsWithCountSlice[i].Count > wordsWithCountSlice[j].Count
	})

	// готовим результат
	if len(wordsWithCountSlice) >= maxResults {
		wordsWithCountSlice = wordsWithCountSlice[0:maxResults]
	}

	result := make([]string, 0, maxResults)
	for _, counter := range wordsWithCountSlice {
		result = append(result, counter.Word)
	}

	return result
}
