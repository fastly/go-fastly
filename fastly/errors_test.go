package fastly

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/jsonapi"
)

func TestNewHTTPError(t *testing.T) {
	t.Parallel()

	t.Run("legacy", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusNotFound,
			Body: io.NopCloser(bytes.NewBufferString(
				`{"msg": "hello", "detail": "nope"}`)),
		}
		e := NewHTTPError(resp)

		if e.StatusCode != http.StatusNotFound {
			t.Errorf("bad status code: %d", e.StatusCode)
		}

		expected := strings.TrimSpace(`
404 - Not Found:

    Title:  hello
    Detail: nope
`)
		if e.Error() != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", e.Error(), expected)
		}
		if e.String() != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", e.String(), expected)
		}

		if !e.IsNotFound() {
			t.Error("not not found")
		}
	})

	t.Run("jsonapi", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header(map[string][]string{"Content-Type": {jsonapi.MediaType}}),
			Body: io.NopCloser(bytes.NewBufferString(
				`{"errors":[{"id":"abc123", "title":"Not found", "detail":"That resource does not exist"}]}`)),
		}
		e := NewHTTPError(resp)

		if e.StatusCode != http.StatusNotFound {
			t.Errorf("expected %d to be %d", e.StatusCode, http.StatusNotFound)
		}

		expected := strings.TrimSpace(`
404 - Not Found:

    ID:     abc123
    Title:  Not found
    Detail: That resource does not exist
`)
		if e.Error() != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", e.Error(), expected)
		}
		if e.String() != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", e.String(), expected)
		}

		if !e.IsNotFound() {
			t.Error("not not found")
		}
	})

	t.Run("problem detail", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header(map[string][]string{"Content-Type": {"application/problem+json"}}),
			Body: io.NopCloser(bytes.NewBufferString(
				`{"title": "Error", "detail": "this was an error", "status": 404}`,
			)),
		}
		e := NewHTTPError(resp)

		if e.StatusCode != http.StatusNotFound {
			t.Errorf("expected %d to be %d", e.StatusCode, http.StatusNotFound)
		}

		expected := strings.TrimSpace(`
404 - Not Found:

    Title:  Error
    Detail: this was an error
`)
		if e.Error() != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", e.Error(), expected)
		}
		if e.String() != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", e.String(), expected)
		}

		if !e.IsNotFound() {
			t.Error("not not found")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		contentTypes := []string{
			jsonapi.MediaType,
			"application/problem+json",
			"default",
		}

		for _, ct := range contentTypes {
			resp := &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     http.Header(map[string][]string{"Content-Type": {ct}}),
				Body:       io.NopCloser(bytes.NewBufferString(`THIS IS NOT JSON`)),
			}
			e := NewHTTPError(resp)

			if e.StatusCode != http.StatusNotFound {
				t.Errorf("expected %d to be %d", e.StatusCode, http.StatusNotFound)
			}

			expected := strings.TrimSpace(`
404 - Not Found:

    Title:  Undefined error
    Detail: THIS IS NOT JSON
`)
			if e.Error() != expected {
				t.Errorf("Content-Type: %q: expected \n\n%q\n\n to be \n\n%q\n\n", ct, e.Error(), expected)
			}
			if e.String() != expected {
				t.Errorf("Content-Type: %q: expected \n\n%q\n\n to be \n\n%q\n\n", ct, e.String(), expected)
			}

			if !e.IsNotFound() {
				t.Error("not not found")
			}
		}
	})
}
