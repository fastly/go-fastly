package fastly

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v13/fastly/impersonation"
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
