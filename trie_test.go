package main

import (
	"reflect"
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
	trie.ComputeSortedURLs()

}

func TestMakeTrieNode(t *testing.T) {
	ptn := MakeTrieNode()
	if ptn == nil {
		t.Errorf("MakeTrieNode() returns nil")
	}
	if ptn.logCounts == nil {
		t.Errorf("MakeTrieNode() should initialize logCounts")
	}
	if ptn.childs == nil {
		t.Errorf("MakeTrieNode() should initialize childs")
	}
	if ptn.sortedUrls != nil {
		t.Errorf("MakeTrieNode() should not initialize childs")
	}
}

func TestMakeTrie(t *testing.T) {
	if MakeTrie().rootNode == nil {
		t.Errorf("MakeTrie() should create a root node")
	}
}

func TestTrieNode_getOrMake(t *testing.T) {
	trie = MakeTrie()
	if len(trie.rootNode.childs) != 0 {
		t.Errorf("MakeTrie() should not create child nodes")
	}
	pTrieNode := trie.rootNode.getOrMake(0)
	if pTrieNode == nil {
		t.Errorf("getOrMake() should return a new node, got nil")
	}
	gotNode, foundNode := trie.rootNode.childs[0]
	if !foundNode || !reflect.DeepEqual(gotNode, pTrieNode) {
		t.Errorf("getOrMake() should properly insert a new node in to childs")
	}
	if len(trie.rootNode.childs) != 1 || trie.rootNode.getOrMake(0) != pTrieNode {
		t.Errorf("getOrMake() should not reinsert a node on the same position")
	}
	var empty *TrieNode
	if empty.getOrMake(1) != nil {
		t.Errorf("getOrMake() should not create nodes on empty recievers")
	}
}

func TestTrieNode_Get(t *testing.T) {
	type fields struct {
		logCounts  stringIntMap
		childs     map[int]*TrieNode
		sortedUrls *[]urlCountPair
	}
	type args struct {
		c []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *TrieNode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pTrieNode := &TrieNode{
				logCounts:  tt.fields.logCounts,
				childs:     tt.fields.childs,
				sortedUrls: tt.fields.sortedUrls,
			}
			if got := pTrieNode.Get(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrieNode.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
