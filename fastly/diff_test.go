package fastly

import (
	"errors"
	"testing"
)

func TestClient_Diff(t *testing.T) {
	t.Parallel()

	var err error
	var tv1 *Version
	Record(t, "diff/version_1", func(c *Client) {
		tv1 = testVersion(t, c)
	})

	var tv2 *Version
	Record(t, "diff/version_2", func(c *Client) {
		tv2 = testVersion(t, c)
	})

	// Diff should be empty
	var d *Diff
	Record(t, "diff/get", func(c *Client) {
		d, err = c.GetDiff(&GetDiffInput{
			ServiceID: TestDeliveryServiceID,
			From:      *tv1.Number,
			To:        *tv2.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a diff
	Record(t, "diff/create_backend", func(c *Client) {
		_, err = c.CreateBackend(&CreateBackendInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv2.Number,
			Name:           ToPointer("test-backend"),
			Address:        ToPointer("integ-test.go-fastly.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure we delete the backend we just created
	defer func() {
		Record(t, "diff/cleanup", func(c *Client) {
			_ = c.DeleteBackend(&DeleteBackendInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv2.Number,
				Name:           "test-backend",
			})
		})
	}()

	// Diff should mot be empty
	Record(t, "diff/get_again", func(c *Client) {
		d, err = c.GetDiff(&GetDiffInput{
			ServiceID: TestDeliveryServiceID,
			From:      *tv1.Number,
			To:        *tv2.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Diff) == 0 {
		t.Errorf("bad diff: %s", d.Diff)
	}
}

func TestClient_Diff_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetDiff(&GetDiffInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDiff(&GetDiffInput{
		ServiceID: "foo",
		From:      0,
	})
	if !errors.Is(err, ErrMissingFrom) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDiff(&GetDiffInput{
		ServiceID: "foo",
		From:      1,
		To:        0,
	})
	if !errors.Is(err, ErrMissingTo) {
		t.Errorf("bad error: %s", err)
	}
}
