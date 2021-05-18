package ipset

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/dghubble/trie"
	"github.com/miekg/dns"
	"strings"
)

// IPSet implements the plugin interface.
type IPSet struct {
	Next     plugin.Handler
	listName []string

	// 匹配的记录
	domains *trie.PathTrie
}

// DotSegment trie path segment
func DotSegment(path string, start int) (segment string, next int) {
	if len(path) == 0 || start < 0 || start > len(path)-1 {
		return "", -1
	}
	end := strings.IndexRune(path[start+1:], '.') // next '/' after 0th rune
	if end == -1 {
		return path[start:], -1
	}
	return path[start : start+end+1], start + end + 1
}

// New create new ipset
func New() *IPSet {
	c := &trie.PathTrieConfig{Segmenter: DotSegment}
	f := &IPSet{domains: trie.NewPathTrieWithConfig(c)}

	return f
}

// ServeDNS implements the plugin.Handler interface.
func (p IPSet) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	wr := NewResponseReverter(w, r, &p)

	if len(p.listName) == 0 {
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, wr, r)
}

// Name implements the Handler interface.
func (p IPSet) Name() string { return "p" }
