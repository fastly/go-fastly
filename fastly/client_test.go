package fastly

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly/impersonation"
)

func TestClient_RawRequest(t *testing.T) {
	validAPIHosts := []string{
		DefaultEndpoint,
		"https://api.fastly.com/",
	}
	purgeAPIPaths := []string{
		"/service/myservice/purge/",
		"service/myservice/purge/",
	}
	cacheKeys := []string{
		"/",
		"text//text",
		"$-_.+!*'(),,;/?:@=&\"<>#%{}|\\^~[]`",
	}
	c := &Client{}
	for _, h := range validAPIHosts {
		var err error
		c.url, err = url.Parse(h)
		if err != nil {
			t.Fatalf("Unable to parse url %s: %s\n", h, err)
		}
		for _, p := range purgeAPIPaths {
			for _, k := range cacheKeys {
				r, err := c.RawRequest(context.TODO(), http.MethodGet, p+url.PathEscape(k), CreateRequestOptions())
				// Cannot test results for success if we get an error
				if err != nil {
					t.Fatal("Could not make RawRequest for ", h, p, k)
				}
				t.Log("Encoded path returned: ", r.URL.EscapedPath())
				pk := p + url.PathEscape(k)
				if p[0] != '/' {
					pk = "/" + pk
				}
				t.Log("Encoded path expected: ", pk)
				// Insure we don't get a path starting with an extra slash
				// e.g. //service/myservice/purge/mykey
				if r.URL.Path[1] == '/' {
					t.Fatalf("Host and APIPath were joined incorrectly. Got: %s\n", r.URL.Path)
				}
				// Insure the encoded path isn't altered
				if !strings.Contains(r.URL.EscapedPath(), p+url.PathEscape(k)) {
					t.Fatalf("RawRequest altered the encoded path. New encoded path: %s, expecting: %s\n", r.URL.EscapedPath(), p+url.PathEscape(k))
				}
			}
		}
	}
}

type impersonationRoundTripper struct {
	rawQuery string
}

func (irt *impersonationRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	irt.rawQuery = req.URL.RawQuery
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func TestClient_Impersonation(t *testing.T) {
	const testCustomerID = "1234ABCD"

	require := require.New(t)
	t.Parallel()

	c, err := NewClient("nokey")
	require.NoError(err)

	irt := &impersonationRoundTripper{}
	ro := CreateRequestOptions()
	c.HTTPClient = &http.Client{Transport: irt}

	_, err = c.Request(impersonation.NewContextForCustomerID(context.TODO(), testCustomerID), http.MethodGet, "/test", ro)
	require.NoError(err)

	require.Equal(impersonation.QueryParam+"="+testCustomerID, irt.rawQuery, "unexpected query parameter")
}

func TestClient_RequestRateLimit(t *testing.T) {
	cases := []struct {
		name          string
		method        string
		status        int
		remaining     string
		reset         string
		wantRemaining int
		wantReset     int64
	}{
		{"success", http.MethodPost, http.StatusOK, "4", "1800000000", 4, 1800000000},
		{"rate limited", http.MethodPost, http.StatusTooManyRequests, "0", "1800000000", 0, 1800000000},
		{"server error", http.MethodDelete, http.StatusInternalServerError, "3", "1800000000", 3, 1800000000},
		{"absent headers", http.MethodPost, http.StatusTooManyRequests, "", "", 10, 1700000000},
		{"invalid headers", http.MethodPost, http.StatusTooManyRequests, "invalid", "invalid", 10, 1700000000},
		{"remaining only", http.MethodPut, http.StatusTooManyRequests, "0", "invalid", 0, 1700000000},
		{"reset only", http.MethodPatch, http.StatusTooManyRequests, "invalid", "1800000000", 10, 1800000000},
		{"get unchanged", http.MethodGet, http.StatusOK, "0", "1800000000", 10, 1700000000},
		{"head unchanged", http.MethodHead, http.StatusOK, "0", "1800000000", 10, 1700000000},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Fastly-RateLimit-Remaining", tc.remaining)
				w.Header().Set("Fastly-RateLimit-Reset", tc.reset)
				w.WriteHeader(tc.status)
			}))
			defer server.Close()
			c, err := NewClientForEndpoint("", server.URL)
			require.NoError(t, err)
			c.remaining = 10
			c.reset = 1700000000
			resp, err := c.Request(context.Background(), tc.method, "/test", CreateRequestOptions())
			require.NotNil(t, resp)
			require.NoError(t, resp.Body.Close())
			if tc.status >= 400 {
				var httpErr *HTTPError
				require.ErrorAs(t, err, &httpErr)
				require.Equal(t, tc.status, httpErr.StatusCode)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.wantRemaining, c.RateLimitRemaining())
			require.Equal(t, time.Unix(tc.wantReset, 0), c.RateLimitReset())
		})
	}
}

func TestClient_RequestRateLimitCanceled(t *testing.T) {
	c, err := NewClientForEndpoint("", "http://127.0.0.1:1")
	require.NoError(t, err)
	c.remaining = 10
	c.reset = 1700000000
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp, err := c.Request(ctx, http.MethodPost, "/test", CreateRequestOptions())
	require.Nil(t, resp)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 10, c.RateLimitRemaining())
	require.Equal(t, time.Unix(1700000000, 0), c.RateLimitReset())
}
