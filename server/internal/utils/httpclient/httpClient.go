// Package httpclient is the factory for outbound *http.Client values. Defaults
// from config: total timeout, optional TLS-insecure flag for self-signed
// internal hosts. The redirect policy is set by targetClient (which needs to
// strip Authorization on cross-host hops, §9.1) — this factory's client has
// the default 10-redirect policy.
package httpclient

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/config"
)

func New(cfg config.Config) *http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: cfg.ExecTimeout,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          50,
		IdleConnTimeout:       90 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: cfg.ExecInsecureTLS, // #nosec G402 — opt-in via TRYIT_EXEC_INSECURE_TLS
		},
	}
	return &http.Client{
		Transport: tr,
		Timeout:   cfg.ExecTimeout,
	}
}
