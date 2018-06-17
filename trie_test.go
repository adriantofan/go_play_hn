package main

import (
	"testing"
	"time"
)

// import "testing"
func makeTestTrie() Trie {
	t := MakeTrie()
	t.AddLog(time.Date(2005, time.January, 1, 1, 3, 15, 0, time.UTC), "url1")
	return t
}

func TestTrie_ComputeTops(t *testing.T) {
	trie := makeTestTrie()
	trie.ComputeTops()

	

}
