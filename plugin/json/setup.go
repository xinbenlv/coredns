package json

import (
	"net/http"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	plugin.Register("json", setup)
}

func setup(c *caddy.Controller) error {
	j := JSON{
		Client: &http.Client{Timeout: 5 * time.Second},
	}

	for c.Next() {
		args := c.RemainingArgs()
		if len(args) != 1 {
			return plugin.Error("json", c.ArgErr())
		}
		j.URL = args[0]
		log.Infof("Configured endpoint: %s", j.URL)

		for c.NextBlock() {
			switch c.Val() {
			case "dnssec":
				j.DNSSEC = true
				log.Info("Enabled DNSSEC support")
			default:
				return plugin.Error("json", c.ArgErr())
			}
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		j.Next = next
		return j
	})

	return nil
} 