package json

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestJSONPlugin(t *testing.T) {
	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"RCODE": 0,
			"AD": true,
			"Answer": [{
				"name": "example.com.",
				"type": 1,
				"TTL": 3600,
				"data": "192.168.1.1"
			}],
			"Question": [{
				"name": "example.com.",
				"type": 1
			}]
		}`
		w.Write([]byte(response))
	}))
	defer ts.Close()

	// Create plugin instance
	j := JSON{
		Client: &http.Client{},
		URL:    ts.URL,
	}

	// Create test context
	ctx := context.TODO()
	rec := dnstest.NewRecorder(&test.ResponseWriter{})

	// Test query
	m := new(dns.Msg)
	m.SetQuestion("example.com.", dns.TypeA)
	
	rcode, err := j.ServeDNS(ctx, rec, m)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rcode != dns.RcodeSuccess {
		t.Errorf("Expected success rcode, got %d", rcode)
	}

	if len(rec.Msg.Answer) != 1 {
		t.Fatalf("Expected 1 answer, got %d", len(rec.Msg.Answer))
	}
} 