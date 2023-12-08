package fastly

import (
	"testing"
)

func TestClient_HTTP3(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "http3/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Enable HTTP3
	var h *HTTP3
	record(t, "http3/enable", func(c *Client) {
		h, err = c.EnableHTTP3(&EnableHTTP3Input{
			FeatureRevision: ToPointer(1),
			ServiceID:       testServiceID,
			ServiceVersion:  *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *h.FeatureRevision != 1 {
		t.Errorf("bad feature_revision: %d", *h.FeatureRevision)
	}

	// Get HTTP3 status
	var gh *HTTP3
	record(t, "http3/get", func(c *Client) {
		gh, err = c.GetHTTP3(&GetHTTP3Input{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *gh.FeatureRevision != 1 {
		t.Errorf("bad feature_revision: %d", *gh.FeatureRevision)
	}

	// Disable HTTP3
	record(t, "http3/disable", func(c *Client) {
		err = c.DisableHTTP3(&DisableHTTP3Input{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get HTTP3 status again to check disabled
	record(t, "http3/get-disabled", func(c *Client) {
		gh, err = c.GetHTTP3(&GetHTTP3Input{
			ServiceID:      testServiceID,
			ServiceVersion: *tv.Number,
		})
	})

	// The API returns a 404 if HTTP3 is not enabled.
	// The API client returns an error if a non-2xx is returned from the API.
	if err == nil {
		t.Fatal("expected a 404 from the API but got a 2xx")
	}
}

func TestClient_GetHTTP3_validation(t *testing.T) {
	var err error

	_, err = testClient.GetHTTP3(&GetHTTP3Input{
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHTTP3(&GetHTTP3Input{
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHTTP3_validation(t *testing.T) {
	var err error
	_, err = testClient.EnableHTTP3(&EnableHTTP3Input{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.EnableHTTP3(&EnableHTTP3Input{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHTTP3_validation(t *testing.T) {
	var err error

	err = testClient.DisableHTTP3(&DisableHTTP3Input{
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DisableHTTP3(&DisableHTTP3Input{
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
