package leetcode_client

import "net/http"

type userAgentTransport struct {
	wrapped http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	return t.wrapped.RoundTrip(req)
}
