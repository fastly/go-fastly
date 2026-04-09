package fastly

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestClient_GetLoggingEndpointErrors_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetLoggingEndpointErrors(context.TODO(), &LoggingEndpointErrorsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetLoggingEndpointErrors(t *testing.T) {
	t.Parallel()

	var err error
	var result *LoggingEndpointErrorsResponse

	// Get logging endpoint errors
	Record(t, "observability_endpoint_error_stream/get", func(c *Client) {
		result, err = c.GetLoggingEndpointErrors(context.TODO(), &LoggingEndpointErrorsInput{
			ServiceID: TestDeliveryServiceID,
			// Timestamps will need to be updated here if you wish to record the API response
			// body. Streamed errors are only maintained for a given period of time.
			From: ToPointer(uint64(1775741900)),
			To:   ToPointer(uint64(1775741920)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Verify pagination links are extracted from Link header
	if result.NextFrom == "" {
		t.Error("expected NextFrom to be populated from Link header")
	}
	if result.PrevFrom == "" {
		t.Error("expected PrevFrom to be populated from Link header")
	}
}

func TestClient_GetLoggingEndpointErrors_with_filters(t *testing.T) {
	t.Parallel()

	var err error
	var result *LoggingEndpointErrorsResponse

	// Get logging endpoint errors with filters
	Record(t, "observability_endpoint_error_stream/get_with_filters", func(c *Client) {
		result, err = c.GetLoggingEndpointErrors(context.TODO(), &LoggingEndpointErrorsInput{
			ServiceID: TestDeliveryServiceID,
			// Timestamps will need to be updated here if you wish to record the API response
			// body. Streamed errors are only maintained for a given period of time.
			From:   ToPointer(uint64(1775741900)),
			To:     ToPointer(uint64(1775741920)),
			Filter: []string{"Broken Log"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Verify pagination links are extracted from Link header
	if result.NextFrom == "" {
		t.Error("expected NextFrom to be populated from Link header")
	}
	if result.PrevFrom == "" {
		t.Error("expected PrevFrom to be populated from Link header")
	}
}

func TestParseLinkHeader(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		wantNext string
		wantPrev string
	}{
		{
			name:     "both next and prev links",
			header:   `</observability/service/kKJb5bOFI47uHeBVluGfX1/logging/errors%3Ffrom=1775741910>; rel="next", </observability/service/kKJb5bOFI47uHeBVluGfX1/logging/errors%3Ffrom=1775741890>; rel="prev"`,
			wantNext: "1775741910",
			wantPrev: "1775741890",
		},
		{
			name:     "only next link",
			header:   `</observability/service/test/logging/errors?from=123>; rel="next"`,
			wantNext: "123",
			wantPrev: "",
		},
		{
			name:     "empty header",
			header:   "",
			wantNext: "",
			wantPrev: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP response with the Link header
			resp := &http.Response{
				Header: http.Header{},
			}
			if tt.header != "" {
				resp.Header.Set("Link", tt.header)
			}

			gotNext, gotPrev := parseLinkHeader(resp)
			if gotNext != tt.wantNext {
				t.Errorf("parseLinkHeader() gotNext = %v, want %v", gotNext, tt.wantNext)
			}
			if gotPrev != tt.wantPrev {
				t.Errorf("parseLinkHeader() gotPrev = %v, want %v", gotPrev, tt.wantPrev)
			}
		})
	}
}
