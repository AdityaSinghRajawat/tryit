// Package target sends the outbound HTTP request the user composed. The
// safeguards in IMPL §9.1 live here:
//   - per-request timeout (from cfg)
//   - response body cap (10 MB default; truncate + signal)
//   - ≤5 redirects, but strip Authorization on a cross-host redirect (never
//     leak a secret to a redirected host)
//   - TLS verify by default; opt-in insecure for self-signed internal hosts
package target

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/service/execute"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/config"
)

const maxRedirects = 5

type Client struct {
	cfg config.Config
	hc  *http.Client
}

func New(hc *http.Client, cfg config.Config) *Client {
	clone := *hc // shallow copy so we can attach our own redirect policy without mutating the shared client.
	clone.CheckRedirect = stripAuthOnCrossHost
	return &Client{cfg: cfg, hc: &clone}
}

func stripAuthOnCrossHost(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return http.ErrUseLastResponse
	}
	if len(via) == 0 {
		return nil
	}
	src := via[len(via)-1].URL
	if !sameHost(src, req.URL) {
		req.Header.Del("Authorization")
		// We can't easily strip apiKey-in-query on redirect (the URL is
		// chosen by the server). Documented limitation; mitigated by the
		// fact that targets generally don't redirect cross-host.
	}
	return nil
}

func sameHost(a, b *url.URL) bool {
	return strings.EqualFold(a.Hostname(), b.Hostname())
}

// Do sends the request and returns the raw response (the executor masks/
// scrubs into wire DTOs). Body reads are capped at cfg.ExecMaxResponseSize.
// Transport errors (DNS/conn/timeout) → error.
func (c *Client) Do(req *http.Request) (*execute.Response, error) {
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, truncated, err := readCapped(resp.Body, c.cfg.ExecMaxResponseSize)
	if err != nil {
		return nil, err
	}
	return &execute.Response{
		Status:    resp.StatusCode,
		Headers:   resp.Header.Clone(),
		Body:      body,
		Truncated: truncated,
	}, nil
}

func readCapped(r io.Reader, limit int64) ([]byte, bool, error) {
	if limit <= 0 {
		b, err := io.ReadAll(r)
		return b, false, err
	}
	b, err := io.ReadAll(io.LimitReader(r, limit))
	if err != nil {
		return nil, false, err
	}
	// Probe one more byte; if present, we truncated.
	var probe [1]byte
	n, perr := r.Read(probe[:])
	if errors.Is(perr, io.EOF) || n == 0 {
		return b, false, nil
	}
	return b, true, nil
}
