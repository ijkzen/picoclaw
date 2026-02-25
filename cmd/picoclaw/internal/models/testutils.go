package models

import (
	"net"
	"net/url"
	"strings"

	"github.com/sipeed/picoclaw/pkg/config"
)

// testModelReachable performs a minimal connectivity check. For HTTP-based
// providers it will attempt a TCP connection to the API base host. This is a
// pragmatic lightweight check used before saving the config.
func testModelReachable(m config.ModelConfig) bool {
	host := ""
	if m.APIBase != "" {
		u, err := url.Parse(m.APIBase)
		if err == nil {
			host = u.Host
		}
	}
	if host == "" {
		// try to infer host from model protocol defaults — fallback to openrouter
		host = "openrouter.ai:443"
	}

	// ensure host has port
	if !strings.Contains(host, ":") {
		host = host + ":443"
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
