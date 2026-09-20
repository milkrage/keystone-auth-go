package keystoneauth

import (
	"fmt"
	"net/http"
)

type Transport struct {
	// Authenticator provides Keystone tokens.
	// Must not be nil.
	Authenticator *Authenticator

	// Base is the base RoundTripper used to make HTTP requests.
	// If nil, http.DefaultTransport is used.
	Base http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response

	for range 2 {
		ctx := req.Context()

		token, err := t.Authenticator.GetToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get keystone token: %w", err)
		}

		clone := req.Clone(ctx)
		clone.Header.Set("X-Auth-Token", token)

		if req.Body != nil {
			if req.GetBody == nil {
				return nil, fmt.Errorf("unable to copy request body: req.GetBody is nil")
			}

			body, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("failed to get request body: %w", err)
			}
			clone.Body = body
		}

		resp, err = t.base().RoundTrip(clone)
		if err != nil {
			return resp, err
		}

		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}

		resp.Body.Close()
		t.Authenticator.InvalidateToken(token)
	}

	return resp, nil
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}

	return http.DefaultTransport
}
