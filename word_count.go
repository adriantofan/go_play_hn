package main

import (
	"strings"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	counts := make(map[string]int)
	for _, word := range words {
		if count, found := counts[word]; found == true {
			counts[word] = count + 1
		} else {
			counts[word] = 1
		}

	}
	return counts
}

// func main() {
// 	wc.Test(WordCount)
// }
