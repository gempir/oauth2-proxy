package requests

import (
	"log"
	"net/http"
	"os"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/version"
)

type verboseTransport struct {
	next      http.RoundTripper
	userAgent string
	logger    *log.Logger
}

func (t *verboseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log request details
	t.logger.Printf("Sending request: Method=%s, URL=%s, UserAgent=%s",
		req.Method, req.URL.String(), req.Header.Get("User-Agent"))

	// Clone the request to avoid modifying the original
	r := req.Clone(req.Context())
	setDefaultUserAgent(r.Header, t.userAgent)

	// Perform the request
	resp, err := t.next.RoundTrip(r)

	// Log any errors or response details
	if err != nil {
		t.logger.Printf("Request failed: Method=%s, URL=%s, Error=%v",
			r.Method, r.URL.String(), err)
		return nil, err
	}

	if resp.StatusCode >= 400 {
		t.logger.Printf("HTTP Error: Status=%d, Method=%s, URL=%s",
			resp.StatusCode, r.Method, r.URL.String())
	}

	return resp, nil
}

var DefaultHTTPClient = &http.Client{Transport: &verboseTransport{
	next:      DefaultTransport,
	userAgent: "oauth2-proxy/" + version.VERSION,
	logger:    log.New(os.Stderr, "[HTTP Client] ", log.LstdFlags),
}}

var DefaultTransport = http.DefaultTransport

func setDefaultUserAgent(header http.Header, userAgent string) {
	if header != nil && len(header.Values("User-Agent")) == 0 {
		header.Set("User-Agent", userAgent)
	}
}
