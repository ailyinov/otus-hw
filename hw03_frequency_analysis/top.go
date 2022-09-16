package hw03frequencyanalysis

import (
	"bufio"
	"math"
	"sort"
	"strings"
)

const topCount = 10

type wordCount struct {
	word  string
	count int
}

type wordsCountList []wordCount

func Top10(in string) (out []string) {
	occurrences := make(map[string]int)
	r := strings.NewReader(in)
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)
	for s.Scan() {
		key := strings.ToLower(strings.Trim(s.Text(), "!.,-?:\"()';"))
		if key != "" {
			occurrences[key]++
		}
	}
	var wordsCount wordsCountList
	for w, c := range occurrences {
		wordsCount = append(wordsCount, wordCount{w, c})
	}
	sort.Slice(wordsCount, sortByCount(wordsCount))

	for _, wc := range wordsCount[:int(math.Min(float64(len(wordsCount)), topCount))] {
		out = append(out, wc.word)
	}

	return out
}

func sortByCount(wcl wordsCountList) func(i int, j int) bool {
	return func(i, j int) bool {
		if wcl[i].count == wcl[j].count {
			return wcl[i].word < wcl[j].word
		}
		return wcl[i].count > wcl[j].count
	}
}
