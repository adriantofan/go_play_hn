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
	trie.ComputeSortedURLs()

}

func TestEqualTrieNode(t *testing.T) {
	type args struct {
		t1 *TrieNode
		t2 *TrieNode
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"both nil",
			args{nil, nil},
			true,
		},
		{
			"both nil",
			args{&TrieNode{}, nil},
			false,
		},
		{
			"both nil",
			args{nil, &TrieNode{}},
			false,
		},
		{
			"diffrent childs length",
			args{&TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: nil},
				nil,
			}, &TrieNode{
				stringIntMap{},
				map[int]*TrieNode{},
				nil,
			}},
			false,
		},
		{
			"diffrent logCounts length",
			args{&TrieNode{
				map[string]int{"2": 1},
				map[int]*TrieNode{},
				nil,
			}, &TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				nil,
			}},
			false,
		},
		{
			"diffrent logCounts count for the same url",
			args{&TrieNode{
				map[string]int{"2": 1},
				map[int]*TrieNode{},
				nil,
			}, &TrieNode{
				map[string]int{"2": 2},
				map[int]*TrieNode{},
				nil,
			}},
			false,
		},
		{
			"logCounts url not present ",
			args{&TrieNode{
				map[string]int{"2": 1},
				map[int]*TrieNode{},
				nil,
			}, &TrieNode{
				map[string]int{"1": 1},
				map[int]*TrieNode{},
				nil,
			}},
			false,
		},
		{
			"diffrent left urlCountPair nil right non nil",
			args{&TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				nil,
			}, &TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				&[]urlCountPair{{"1", 1}},
			}},
			false,
		},
		{
			"diffrent left urlCountPair non nil right nil",
			args{&TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				&[]urlCountPair{{"1", 1}},
			}, &TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				nil,
			}},
			false,
		},
		{
			"diffrent urlCountPair length",
			args{&TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				&[]urlCountPair{{"1", 1}},
			}, &TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				&[]urlCountPair{},
			}},
			false,
		},
		{
			"urlCountPair same lenght diffrent value",
			args{&TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				&[]urlCountPair{{"1", 1}},
			}, &TrieNode{
				map[string]int{},
				map[int]*TrieNode{},
				&[]urlCountPair{{"1", 2}},
			}},
			false,
		},
		{
			"childs not equal",
			args{&TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: &TrieNode{
					stringIntMap{},
					map[int]*TrieNode{1: nil},
					nil,
				}},
				nil,
			}, &TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: &TrieNode{
					stringIntMap{},
					map[int]*TrieNode{},
					nil,
				}},
				nil,
			}},
			false,
		},
		{
			"child left nil right non nil",
			args{&TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: nil},
				nil,
			}, &TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: &TrieNode{}},
				nil,
			}},
			false,
		},
		{
			"child left non nil right nil",
			args{&TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: &TrieNode{}},
				nil,
			}, &TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: nil},
				nil,
			}},
			false,
		},
		{
			"child left not equal right",
			args{&TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: &TrieNode{}},
				nil,
			}, &TrieNode{
				stringIntMap{},
				map[int]*TrieNode{1: nil},
				nil,
			}},
			false,
		},
		{
			"equal",
			args{&TrieNode{
				map[string]int{"2": 1},
				map[int]*TrieNode{1: &TrieNode{}},
				&[]urlCountPair{{"1", 2}},
			}, &TrieNode{
				map[string]int{"2": 1},
				map[int]*TrieNode{1: &TrieNode{}},
				&[]urlCountPair{{"1", 2}},
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.t1.Equal(tt.args.t2); got != tt.want {
				t.Errorf("EqualTrieNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
