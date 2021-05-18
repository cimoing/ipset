package ipset

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// IPSet implements the plugin interface.
type IPSet struct {
	Next     plugin.Handler
	listName []string
}

// ServeDNS implements the plugin.Handler interface.
func (p IPSet) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	wr := NewResponseReverter(w, r, p.listName)

	if len(p.listName) == 0 {
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, wr, r)
}

// Name implements the Handler interface.
func (p IPSet) Name() string { return "p" }
