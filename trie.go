package ipset

import "strings"

type StringSegment func(key string, start int) (segment string, nextIndex int)

type DotSegmentIterator struct {
	current int
	value   string
	revert  bool
}

func NewDotSegmentRevertIterator(value string) *DotSegmentIterator {
	return &DotSegmentIterator{value: value, revert: true, current: len(value)}
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
	key      string
	value    int
	children map[string]*Trie
}

func NewTrie(key string) *Trie {
	return &Trie{key: key, value: -1, children: make(map[string]*Trie)}
}

func (trie *Trie) Get(key string) int {
	iterator := NewDotSegmentRevertIterator(key)

	node := trie
	current := iterator.current
	next := -1

	for {
		current = iterator.current
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
		current = iterator.current
		next = iterator.Next()

		// get to end
		if next == -1 {
			node.value = value
			break
		}

		if nextNode, ok := node.children[iterator.value[next:current]]; ok {
			node = nextNode
		} else {
			tmp := NewTrie(iterator.value[next:current])
			node.children[iterator.value[next:current]] = tmp
			node = tmp
		}
	}
}
