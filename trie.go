package main

// Trie datastructure for log processing
type Trie struct {
	rootNode *TrieNode
}

// TrieNode implementation for storing log lines organized by time
// could log count be generic ? so we might keep a seccond trie with the sorted logs and dropping the one having the counts
type TrieNode struct {
	logCounts  stringIntMap
	childs     map[int]*TrieNode
	sortedUrls *[]QueryCountPair
}

// QueryCountPair used internaly and in kind'of hackish way to display the json
type QueryCountPair struct {
	Query string `json:"query"`
	Count int    `json:"count"`
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

func (pTrieNode *TrieNode) getOrMake(component int) (child *TrieNode) {
	if pTrieNode == nil {
		return
	}
	child, found := pTrieNode.childs[component]
	if found {
		return
	}
	//BUG(atn) remove reference to MakeTrie wich creates a specific trie with urlCount backed by a hash map
	child = MakeTrieNode()
	pTrieNode.childs[component] = child
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
func (pTrieNode *TrieNode) Get(c []int) *TrieNode {
	if pTrieNode != nil {
		// base case
		if len(c) == 0 {
			return pTrieNode
		}
		child, hasChild := pTrieNode.childs[c[0]]
		if hasChild {
			// child is a leaf
			if len(child.childs) == 0 {
				return child
			}
			// recursive case
			return child.Get(c[1:])
		}
	}
	// not found
	return nil
}

// Visit calls handler eagerly for all reachable nodes starting with self
func (pTrieNode *TrieNode) Visit(handler func(*TrieNode)) {
	if pTrieNode == nil {
		return
	}
	handler(pTrieNode)
	for _, pChild := range pTrieNode.childs {
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
