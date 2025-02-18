package json

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// log is the plugin logger (single declaration)
var log = clog.NewWithPlugin("json")

type JSON struct {
	Next   plugin.Handler
	Client *http.Client
	URL    string
	DNSSEC bool // TODO: implement dnssec signing
}

type DNSResponse struct {
	RCODE   int            `json:"RCODE"`
	AD      bool           `json:"AD"`
	Answer  []DNSAnswer    `json:"Answer"`
	Question []DNSQuestion `json:"Question"`
}

type DNSAnswer struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
	TTL  uint32 `json:"TTL"`
	Data string `json:"data"`
}

type DNSQuestion struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
}

func (j JSON) Name() string { return "json" }

func (j JSON) ServeDNS(ctx context.Context, w dns.ResponseWriter, m *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: m}
	qname := state.Name()
	
	log.Debugf("Received query for %s (type %d)", qname, state.QType())
	log.Debugf("Constructing API URL: %s with qname: %s", j.URL, qname)

	// Build REST API URL with query name
	url := fmt.Sprintf("%s%s", j.URL, qname)
	log.Debugf("Building URL: %s", url)
	// Create HTTP request with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Debugf("Creating HTTP request with context: %v", ctx)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return dns.RcodeServerFailure, err
	}

	// Execute HTTP request
	resp, err := j.Client.Do(req)
	if err != nil {
		log.Errorf("HTTP request failed: %v", err)
		return dns.RcodeServerFailure, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Unexpected status code: %d from %s", resp.StatusCode, url)
		return dns.RcodeServerFailure, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	// Parse JSON response
	var dnsResp DNSResponse
	if err := json.NewDecoder(resp.Body).Decode(&dnsResp); err != nil {
		log.Errorf("Failed to decode response: %v", err)
		return dns.RcodeServerFailure, err
	}

	log.Debugf("API response - RCODE: %d, Answers: %d", dnsResp.RCODE, len(dnsResp.Answer))

	// Create DNS response message
	reply := new(dns.Msg)
	reply.SetReply(m)
	reply.Rcode = dnsResp.RCODE
	reply.AuthenticatedData = dnsResp.AD

	// Convert JSON answers to DNS RRs
	for _, ans := range dnsResp.Answer {
		rr, err := dns.NewRR(fmt.Sprintf("%s %d %s %s", ans.Name, ans.TTL, dns.Type(ans.Type), ans.Data))
		if err == nil {
			reply.Answer = append(reply.Answer, rr)
		}
	}

	w.WriteMsg(reply)
	return dns.RcodeSuccess, nil
} 