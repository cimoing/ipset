package ipset

import "strings"

type StringSegment func(key string, start int) (segment string, nextIndex int)

// DotSegment trie path segment
func DotSegment(path string, end int) (segment string, next int) {
	if len(path) == 0 || end <= 0 {
		return "", -1
	}
	start := strings.IndexRune(path[end+1:], '.') // next '/' after 0th rune

	if start == -1 {
		return path[end:], -1
	}
	return path[end-start-1 : end], end - start - 1
}

type DotSegmentIterator struct {
	current int
	value   string
	revert  bool
}

func NewDotSegmentRevertIterator(value string) *DotSegmentIterator {
	return &DotSegmentIterator{value: value, revert: true, current: len(value) - 1}
}

func (iterator *DotSegmentIterator) Next() int {
	n := iterator.NextIdx()
	if n == -1 {
		return -1
	}
	iterator.current = n

	return n
}

func (iterator *DotSegmentIterator) HasNext() bool {
	return iterator.NextIdx() != -1
}

func (iterator *DotSegmentIterator) NextIdx() int {
	if iterator.current == -1 {
		return -1
	}
	if iterator.revert {
		if iterator.current <= 0 {
			return -1
		}
		return strings.LastIndexByte(iterator.value[0:iterator.current], '.')
	} else {
		if iterator.current >= len(iterator.value) {
			return -1
		}
		return strings.IndexRune(iterator.value[iterator.current:], '.')
	}
}

type Trie struct {
	segment  StringSegment
	value    int
	children map[string]*Trie
}

func NewTrie() *Trie {
	return &Trie{value: -1, children: make(map[string]*Trie)}
}

func (trie *Trie) Get(key string) int {
	iterator := NewDotSegmentRevertIterator(key)

	node := trie
	current := iterator.current
	next := -1

	for {
		next = iterator.Next()

		// get to end
		if next == -1 {
			break
		}

		if nextNode, ok := node.children[iterator.value[next:current]]; ok {
			node = nextNode
		} else {
			break
		}
	}

	return node.value
}

func (trie *Trie) Put(key string, value int) {
	iterator := NewDotSegmentRevertIterator(key)

	node := trie
	current := iterator.current
	next := -1

	for {
		next = iterator.Next()

		// get to end
		if next == -1 {
			node.value = value
			break
		}

		if nextNode, ok := node.children[iterator.value[next:current]]; ok {
			node = nextNode
		} else {
			node.children[iterator.value[next:current]] = NewTrie()
		}
	}
}
