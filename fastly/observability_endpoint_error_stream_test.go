package fastly

import (
	"context"
	"errors"
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
			From: ToPointer(uint64(1775587245)),
			To:   ToPointer(uint64(1775587545)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
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
			From:   ToPointer(uint64(1775587245)),
			To:     ToPointer(uint64(1775587545)),
			Filter: []string{"Broken Log"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
