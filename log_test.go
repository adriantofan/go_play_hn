package main

import (
	"reflect"
	"testing"
)

func TestTrie_AddLog(t *testing.T) {
	trie := MakeTrie()
	trie.AddLog(ParseTime("2015-08-01 00:03:43"), "term")
	trie.AddLog(ParseTime("2015-08-01 00:03:43"), "term")

	nodeCount := 0
	termCount := 0
	wantedNodeCount := 6
	wantedTermCount := 12
	trie.Visit(func(pTrieNode *TrieNode) {
		nodeCount++
		if count, foundURL := pTrieNode.logCounts["term"]; foundURL {
			termCount = termCount + count
		}
	})
	if nodeCount != wantedNodeCount {
		t.Errorf("AddLog() should add url to all levels got %v want %v", nodeCount, wantedNodeCount)
	}
	if termCount != wantedTermCount {
		t.Errorf("AddLog() total count for 'url' mismatch, got %v want %v", termCount, wantedTermCount)
	}
}

func TestTrie_ComputeSortedURLs(t *testing.T) {
	trie := MakeTrie()
	trie.AddLog(ParseTime("2015-08-01 00:03:43"), "term1")
	trie.AddLog(ParseTime("2015-08-01 00:03:43"), "term2")
	trie.AddLog(ParseTime("2015-08-02 00:03:43"), "term1")
	trie.ComputeSortedURLs()
	trie.Visit(func(pTrieNode *TrieNode) {
		if pTrieNode.logCounts != nil {
			t.Errorf("ComputeSortedURLs should clear logCounts got %v", pTrieNode.logCounts)
		}
	})
	tests := []struct {
		timeArray       []int
		wantSortedTerms []urlCountPair
	}{
		{
			[]int{2015},
			[]urlCountPair{urlCountPair{"term1", 2}, urlCountPair{"term2", 1}},
		},
		{
			[]int{2015, 8},
			[]urlCountPair{urlCountPair{"term1", 2}, urlCountPair{"term2", 1}},
		},
		{
			[]int{2015, 8, 1},
			[]urlCountPair{urlCountPair{"term1", 1}, urlCountPair{"term2", 1}},
		},
		{
			[]int{2015, 8, 1, 0},
			[]urlCountPair{urlCountPair{"term1", 1}, urlCountPair{"term2", 1}},
		},
		{
			[]int{2015, 8, 1, 0, 3},
			[]urlCountPair{urlCountPair{"term1", 1}, urlCountPair{"term2", 1}},
		},
		{
			[]int{2015, 8, 1, 0, 3, 43},
			[]urlCountPair{urlCountPair{"term1", 1}, urlCountPair{"term2", 1}},
		},
	}
	// BUG(adrian) for the same url count the result is not deterministic  [{term1 1} {term2 1}] sortedUrls &[{term2 1} {term1 1}]
	for _, d := range tests {
		gotChild := trie.Get(d.timeArray)
		if gotChild == nil {
			t.Errorf("ComputeSortedURLs got nil for %v", d.timeArray)
			return
		}
		if !reflect.DeepEqual(&d.wantSortedTerms, gotChild.sortedUrls) {
			t.Errorf("ComputeSortedURLs for  gotSortedItems %v sortedUrls %v", d.wantSortedTerms, gotChild.sortedUrls)
		}
	}

}