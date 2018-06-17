package main

import (
	"sort"
	"time"
)

// Trie datastructure for log processing
type Trie struct {
	rootNode *TrieNode
}

// TrieNode implementation for storing log lines organized by time
// could log count be generic ? so we might keep a seccond trie with the sorted logs and dropping the one having the counts
type TrieNode struct {
	logCounts  stringIntMap
	childs     map[int]*TrieNode
	sortedUrls *[]urlCountPair
}

type urlCountPair struct {
	url   string
	count int
}

// MakeTrie creates a Trie initialized with a root node
func MakeTrie() Trie {
	return Trie{MakeTrieNode()}
}

// MakeTrieNode creates a TrieNode with logCount baked by a hash map and returns a reference to it
func MakeTrieNode() *TrieNode {
	t := new(TrieNode)
	t.logCounts = make(stringIntMap)
	t.childs = make(map[int]*TrieNode)
	return t
}

func (n *TrieNode) getOrMake(component int) (child *TrieNode) {
	child, found := n.childs[component]
	if found {
		return
	}
	//BUG(atn) remove reference to MakeTrie wich creates a specific trie with urlCount backed by a hash map
	child = MakeTrieNode()
	n.childs[component] = child
	return
}

// Get returns the note with the specified prefix in trie t or nil if not found
func (t Trie) Get(c []int) *TrieNode {
	if t.rootNode == nil {
		return nil
	}
	return t.rootNode.Get(c)
}

// Get returns the node with the specified prefix in trie t or nil if not found
func (n *TrieNode) Get(c []int) *TrieNode {
	// base case
	if len(c) == 0 {
		return n
	}
	child, hasChild := n.childs[c[0]]
	if hasChild {
		// recursive case
		return child.Get(c[1:])
	}
	// not found
	return nil
}

// Visit calls handler eagerly for all reachable nodes starting with self
func (n *TrieNode) Visit(handler func(*TrieNode)) {
	handler(n)
	for _, pChild := range n.childs {
		pChild.Visit(handler)
	}
}

// Visit calls handler eagerly for all reachable nodes starting from root
func (t Trie) Visit(handler func(*TrieNode)) {
	if t.rootNode == nil {
		return
	}
	t.rootNode.Visit(handler)

}

// a wrapper arrounf a map
type stringIntMap map[string]int

// a simple interface to keep counts a set of strings
type urlCount interface {
	increase(string)
}

// used to implement urlCount interface on a stringIntMap
func (m stringIntMap) increase(url string) {
	count, found := m[url]
	if found {
		m[url] = count + 1
	} else {
		m[url] = 1
	}
}

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

func (t Trie) ComputeTops() {
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
