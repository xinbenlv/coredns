package rest

import (
	"net/http"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	plugin.Register("rest", setup)
}

func setup(c *caddy.Controller) error {
	r := Rest{
		Client: &http.Client{Timeout: 5 * time.Second},
	}

	for c.Next() {
		args := c.RemainingArgs()
		if len(args) != 1 {
			return plugin.Error("rest", c.ArgErr())
		}
		r.URL = args[0]
		log.Infof("Configured endpoint: %s", r.URL)

		for c.NextBlock() {
			switch c.Val() {
			case "dnssec":
				r.DNSSEC = true
				log.Info("Enabled DNSSEC support")
			default:
				return plugin.Error("rest", c.ArgErr())
			}
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		r.Next = next
		return r
	})

	return nil
} 