package leetcode

import "net/http"

type UserAgentTransport struct {
	Wrapped http.RoundTripper
}

func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	return t.Wrapped.RoundTrip(req)
}
