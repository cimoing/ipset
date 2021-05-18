package ipset

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/coredns/caddy"
)

func init() { plugin.Register("ipset", setup) }

func setup(c *caddy.Controller) error {
	var listName []string
	for c.Next() {
		args := c.RemainingArgs()
		listName = append(listName, args...)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return IPSet{Next: next, listName: listName}
	})

	if len(listName) > 0 {
		c.OnStartup(func() error {
			return initLib()
		})

		c.OnShutdown(func() error {
			return shutdownLib()
		})
	}

	return nil
}
