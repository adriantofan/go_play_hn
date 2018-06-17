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
		t.Errorf("MakeTrieNode() should initialize child map")
	}
	if ptn.sortedUrls != nil {
		t.Errorf("MakeTrieNode() should not initialize child map")
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
		t.Errorf("getOrMake() should properly insert a new node in to child array")
	}
	if len(trie.rootNode.childs) != 1 || trie.rootNode.getOrMake(0) != pTrieNode {
		t.Errorf("getOrMake() should not reinsert a node on the same position")
	}
	var empty *TrieNode
	if empty.getOrMake(1) != nil {
		t.Errorf("getOrMake() should not create nodes on empty receivers")
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
		{"nil when child exists and it doesn't contain it",
			&TrieNode{
				make(map[string]int),
				map[int]*TrieNode{1: MakeTrieNode()},
				nil,
			},
			args{[]int{2}},
			nil,
		},
		{"finds it when child exists and it contains it",
			&TrieNode{
				make(map[string]int),
				map[int]*TrieNode{1: empty1},
				nil,
			},
			args{[]int{1}},
			empty1,
		},
		{"finds it when child exists and it contains it deeper",
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

// refactor to go spec
func TestTrieNode_Visit(t *testing.T) {
	var empty *TrieNode
	visited := make([]*TrieNode, 0)
	handler := func(t *TrieNode) {
		visited = append(visited, t)
	}
	empty.Visit(handler)
	if len(visited) != 0 {
		t.Errorf("TrieNode.Visit should not visit nil")
	}
	one := MakeTrieNode()
	one.Visit(handler)
	if len(visited) != 1 {
		t.Errorf("TrieNode.Visit should visit one")
	}
	visitedOne := []*TrieNode{one}
	if !reflect.DeepEqual(visited, visitedOne) {
		t.Errorf("TrieNode.Visit should visit one and call handler with one")
	}
	visited = make([]*TrieNode, 0)
	nested := MakeTrieNode()
	nestedOne := nested.getOrMake(1)
	nestedOneTwo := nestedOne.getOrMake(11)
	nested.Visit(handler)
	visitedNested := []*TrieNode{nested, nestedOne, nestedOneTwo}
	if !reflect.DeepEqual(visited, visitedNested) {
		t.Errorf("TrieNode.Visit should visit nested nodes. got %v , wanting %v", visited, visitedNested)
	}
}

func TestTrie_Visit(t *testing.T) {
	var empty = Trie{}

	visited := make([]*TrieNode, 0)
	handler := func(t *TrieNode) {
		visited = append(visited, t)
	}

	empty.Visit(handler)
	if len(visited) != 0 {
		t.Errorf("TrieNode.Visit should not visit nil")
	}
	one := MakeTrie()
	one.Visit(handler)
	if len(visited) != 1 {
		t.Errorf("Trie.Visit should visit one")
	}
}

func Test_stringIntMap_increase(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		m    stringIntMap
		args args
		want stringIntMap
	}{
		{
			"set when not found",
			map[string]int{},
			args{"one"},
			map[string]int{"one": 1},
		},
		{
			"set when not found",
			map[string]int{"one": 1, "two": 1},
			args{"one"},
			map[string]int{"one": 2, "two": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.increase(tt.args.url)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("TrieNode.increase() -> %v, want %v", tt.m, tt.want)
			}
		})
	}
}
