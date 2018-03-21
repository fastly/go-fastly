package fastly

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/google/jsonapi"
)

func TestNewHTTPError(t *testing.T) {
	t.Parallel()

	t.Run("legacy", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 404,
			Body: ioutil.NopCloser(bytes.NewBufferString(
				`{"msg": "hello", "detail": "nope"}`)),
		}
		e := NewHTTPError(resp)

		if e.StatusCode != 404 {
			t.Errorf("bad status code: %d", e.StatusCode)
		}

		expected := strings.TrimSpace(`
404 - Not Found:

    Title:  hello
    Detail: nope
`)
		if e.Error() != expected {
			t.Errorf("expected \n\n%s\n\n to be \n\n%s\n\n", e.Error(), expected)
		}
		if e.String() != expected {
			t.Errorf("expected \n\n%s\n\n to be \n\n%s\n\n", e.String(), expected)
		}

		if !e.IsNotFound() {
			t.Error("not not found")
		}
	})

	t.Run("jsonapi", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 404,
			Header:     http.Header(map[string][]string{"Content-Type": []string{jsonapi.MediaType}}),
			Body: ioutil.NopCloser(bytes.NewBufferString(
				`{"errors":[{"id":"abc123", "title":"Not found", "detail":"That resource does not exist"}]}`)),
		}
		e := NewHTTPError(resp)

		if e.StatusCode != 404 {
			t.Errorf("expected %d to be %d", e.StatusCode, 404)
		}

		expected := strings.TrimSpace(`
404 - Not Found:

    ID:     abc123
    Title:  Not found
    Detail: That resource does not exist
`)
		if e.Error() != expected {
			t.Errorf("expected \n\n%s\n\n to be \n\n%s\n\n", e.Error(), expected)
		}
		if e.String() != expected {
			t.Errorf("expected \n\n%s\n\n to be \n\n%s\n\n", e.String(), expected)
		}

		if !e.IsNotFound() {
			t.Error("not not found")
		}
	})
}
