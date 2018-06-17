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

func (pTrieNode *TrieNode) getOrMake(component int) (child *TrieNode) {
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
	// base case
	if len(c) == 0 {
		return pTrieNode
	}
	child, hasChild := pTrieNode.childs[c[0]]
	if hasChild {
		// recursive case
		return child.Get(c[1:])
	}
	// not found
	return nil
}

// Visit calls handler eagerly for all reachable nodes starting with self
func (pTrieNode *TrieNode) Visit(handler func(*TrieNode)) {
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

// Equal returns true if p1 is semantiqualy equal to p2
func (p1 urlCountPair) Equal(p2 urlCountPair) bool {
	return p1.url == p2.url && p1.count == p2.count
}

// Equal returns true if the two tries are semanticaly equal
func (t Trie) Equal(trie2 Trie) bool {
	return t.rootNode.Equal(trie2.rootNode)
}

// Equal returns true if pTrieNode is semantiqualy equal to trieNode2
func (pTrieNode *TrieNode) Equal(trieNode2 *TrieNode) bool {
	if pTrieNode == trieNode2 {
		return true
	}
	if (pTrieNode == nil && trieNode2 != nil) || (pTrieNode != nil && trieNode2 == nil) {
		return false
	}
	if pTrieNode == nil && trieNode2 == nil {
		return true
	}
	if pTrieNode.logCounts == nil && trieNode2.logCounts != nil ||
		pTrieNode.logCounts != nil && trieNode2.logCounts == nil {
		return false
	}
	if len(pTrieNode.childs) != len(trieNode2.childs) ||
		len(pTrieNode.childs) != len(trieNode2.childs) {
		return false
	}
	if pTrieNode.sortedUrls == nil && trieNode2.sortedUrls != nil ||
		pTrieNode.sortedUrls != nil && trieNode2.sortedUrls == nil {
		return false
	}

	if pTrieNode.sortedUrls != nil && trieNode2.sortedUrls != nil {
		if len(*pTrieNode.sortedUrls) != len(*trieNode2.sortedUrls) {
			return false
		}
		for i := 0; i < len(*pTrieNode.sortedUrls); i++ {
			if !(*pTrieNode.sortedUrls)[i].Equal((*trieNode2.sortedUrls)[i]) {
				return false
			}
		}
	}
	for logURL1, logCounpTrieNode1 := range pTrieNode.logCounts {
		logCountrieNode2, foundCountrieNode2 := trieNode2.logCounts[logURL1]
		if !foundCountrieNode2 || logCountrieNode2 != logCounpTrieNode1 {
			return false
		}
	}
	for key1, childP1 := range pTrieNode.childs {
		childP2, foundKey2 := trieNode2.childs[key1]
		if !foundKey2 {
			return false
		}
		if childP1 != nil && childP2 == nil ||
			childP1 == nil && childP2 != nil {
			return false
		}
		if childP1 != nil && childP2 != nil {
			if !childP1.Equal(childP2) {
				return false
			}
		}
	}
	return true
}
