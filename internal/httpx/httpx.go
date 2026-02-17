package httpx

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
const DefaultTorProxyURL = "socks5://127.0.0.1:9050"

// Doer lets us accept *http.Client or a test double.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	Timeout     time.Duration
	WithTor     bool
	TorProxyURL string
}

func NewClient(cfg ClientConfig) (*http.Client, error) {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}
	if cfg.TorProxyURL == "" {
		cfg.TorProxyURL = DefaultTorProxyURL
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if cfg.WithTor {
		u, err := url.Parse(cfg.TorProxyURL)
		if err != nil {
			return nil, fmt.Errorf("parse tor proxy url: %w", err)
		}

		dialer, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("create tor dialer: %w", err)
		}

		transport.Proxy = nil
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			// x/net/proxy Dialer doesn't support ctx; best effort.
			return dialer.Dial(network, addr)
		}
	}

	return &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}, nil
}

func NewRequest(ctx context.Context, method, rawURL string, body io.Reader, userAgent string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	return req, nil
}
