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
	type args struct {
		c []int
	}
	empty1 := &TrieNode{
		make(map[string]int),
		make(map[int]*TrieNode),
		nil,
	}
	notEmpty := &TrieNode{
		make(map[string]int),
		map[int]*TrieNode{1: empty1},
		nil,
	}
	tests := []struct {
		name      string
		pTrieNode *TrieNode
		args      args
		want      *TrieNode
	}{
		{"when path is empty return self",
			empty1,
			args{[]int{}},
			empty1,
		},
		{"nil when childs exists and it doesen't contain it",
			&TrieNode{
				make(map[string]int),
				map[int]*TrieNode{1: MakeTrieNode()},
				nil,
			},
			args{[]int{2}},
			nil,
		},
		{"finds it when childs exists and it contains it",
			&TrieNode{
				make(map[string]int),
				map[int]*TrieNode{1: empty1},
				nil,
			},
			args{[]int{1}},
			empty1,
		},
		{"finds it when childs exists and it contains it deeper",
			&TrieNode{
				make(map[string]int),
				map[int]*TrieNode{2: notEmpty},
				nil,
			},
			args{[]int{2, 1}},
			empty1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pTrieNode.Get(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrieNode.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
