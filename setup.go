package database

import (
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
)

func init() {
	plugin.Register(Name, setup)
}

func setup(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		return plugin.Error(Name, c.ArgErr())
	}

	c.OnStartup(func() error {
		once.Do(func() {
			metrics.MustRegister(c, requestCount)
		})
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return DBBackend{Next: next, Debug: true}
	})

	return nil
}
