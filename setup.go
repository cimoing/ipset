package ipset

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/coredns/caddy"
)

func init() { plugin.Register("ipset", setup) }

func setup(c *caddy.Controller) error {
	f, err := parseIPSet(c)

	if err != nil {
		return plugin.Error("ipset", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		f.Next = next
		return f
	})

	if len(f.listName) > 0 {
		c.OnStartup(func() error {
			return initLib()
		})

		c.OnShutdown(func() error {
			return shutdownLib()
		})
	}

	return nil
}

func parseIPSet(c *caddy.Controller) (*IPSet, error) {
	var (
		f   *IPSet
		err error
		i   int
	)

	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++
		f, err = parseStanza(c)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func parseStanza(c *caddy.Controller) (*IPSet, error) {
	f := New()

	args := c.RemainingArgs()
	f.listName = append(f.listName, args...)

	for c.NextBlock() {
		if err := parseBlock(c, f); err != nil {
			return f, err
		}
	}

	return f, nil
}

func parseBlock(c *caddy.Controller, f *IPSet) error {
	val := -1
	switch c.Val() {
	case "include":
		val = 1
	case "exclude":
		val = 0
	}

	domains := c.RemainingArgs()
	if len(domains) == 0 {
		return c.ArgErr()
	}

	for i := 0; i < len(domains); i++ {
		f.domains.Put(domains[i], val)
	}

	return nil
}
