package main

import (
	"sort"
	"time"
)

// AddLog navigates the trie down and adds the url to each date component
func (n *TrieNode) AddLog(components []int, url string) {
	n.logCounts.increase(url)
	// add the url to subsequent levels
	if len(components) != 1 {
		n.getOrMake(components[0]).AddLog(components[1:], url)
	}
}

// AddLog parses time in to a date component array and uses and uses root node add to pass the url allong the trie
func (t Trie) AddLog(date time.Time, url string) {
	components := LogDateComponents(date)
	t.rootNode.AddLog(components[:], url)
}

func (t Trie) ComputeSortedURLs() {
	t.Visit(func(n *TrieNode) {
		sorted := []urlCountPair{}
		if n.logCounts == nil {
			return
		}
		for url, count := range n.logCounts {
			sorted = append(sorted, urlCountPair{url, count})
		}
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].count > sorted[j].count
		})
		n.sortedUrls = &sorted
		n.logCounts = nil
	})
}

// TopNAtDate returns top n urls at the given date where c contains the most significants components of that date
// [Year, Month, Day, Hour, Minute, Seccond]. For example in 2012 c is [2012]; in 2012-12 c is [2012, 12]
func TopNAtDate(t Trie, c []int, n int) []urlCountPair {
	node := t.Get(c)
	if node != nil && node.sortedUrls != nil && len(*node.sortedUrls) > 0 {
		if n > len(*node.sortedUrls) {
			n = len(*node.sortedUrls)
		}
		result := (*node.sortedUrls)[:n]
		return result
	}
	return []urlCountPair{}
}

// Distinct returns the how many distinct urls are at the given date, where c contains the most significants components of the date
// [Year, Month, Day, Hour, Minute, Seccond]. For example in 2012 c is [2012]; in 2012-12 c is [2012, 12]
func Distinct(t Trie, c []int) int {
	node := t.Get(c)
	if node != nil && node.sortedUrls != nil {
		return len(*node.sortedUrls)
	}
	return 0
}