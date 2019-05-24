package fastly

import (
	"net/url"
	"strings"
	"testing"
)

func TestClient_RawRequest(t *testing.T) {
	validAPIHosts := []string{
		"https://api.fastly.com",
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
				r, err := c.RawRequest("GET", p+url.PathEscape(k), nil)
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
				if strings.Index(r.URL.EscapedPath(), p+url.PathEscape(k)) == -1 {
					t.Fatalf("RawRequest altered the encoded path. New encoded path: %s, expecting: %s\n", r.URL.EscapedPath(), p+url.PathEscape(k))
				}
			}
		}
	}
}
