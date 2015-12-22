package fastly

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewHTTPError(t *testing.T) {
	resp := &http.Response{
		StatusCode: 404,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"msg": "hello", "detail": "nope"}`)),
	}
	e := NewHTTPError(resp)

	if e.StatusCode != 404 {
		t.Errorf("bad status code: %d", e.StatusCode)
	}

	if e.Message != "hello" {
		t.Errorf("bad message: %q", e.Message)
	}

	if e.Detail != "nope" {
		t.Errorf("bad detail: %q", e.Detail)
	}

	if e.Error() != "404 - Not Found\nMessage: hello\nDetail: nope" {
		t.Errorf("bad error: %q", e)
	}

	if e.String() != "404 - Not Found\nMessage: hello\nDetail: nope" {
		t.Errorf("bad string: %q", e)
	}

	if !e.IsNotFound() {
		t.Error("not not found")
	}
}
