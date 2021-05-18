package ipset

import "testing"

func TestTriePutGet(t *testing.T) {
	trie := NewTrie("")
	trie.Put(".google.com", 1)
	trie.Put(".dl.google.com", 0)
	trie.Put(".facebook.com", 1)

	if trie.Get("www.google.com") != 1 {
		t.Error("not match")
	}

	if trie.Get("google.com") != 1 {
		t.Error("not match")
	}

	if trie.Get("dl.google.com") != 0 {
		t.Error("not match")
	}

	v := trie.Get("www.foo.com")
	if v != -1 {
		t.Error("should not match", v)
	}
}
