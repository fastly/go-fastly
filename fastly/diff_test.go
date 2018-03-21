package fastly

import "testing"

func TestClient_Diff(t *testing.T) {
	t.Parallel()

	var err error
	var tv1 *Version
	record(t, "diff/version_1", func(c *Client) {
		tv1 = testVersion(t, c)
	})

	var tv2 *Version
	record(t, "diff/version_2", func(c *Client) {
		tv2 = testVersion(t, c)
	})

	// Diff should be empty
	var d *Diff
	record(t, "diff/get", func(c *Client) {
		d, err = c.GetDiff(&GetDiffInput{
			Service: testServiceID,
			From:    tv1.Number,
			To:      tv2.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a diff
	record(t, "diff/create_backend", func(c *Client) {
		_, err = c.CreateBackend(&CreateBackendInput{
			Service: testServiceID,
			Version: tv2.Number,
			Name:    "test-backend",
			Address: "integ-test.go-fastly.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure we delete the backend we just created
	defer func() {
		record(t, "diff/cleanup", func(c *Client) {
			c.DeleteBackend(&DeleteBackendInput{
				Service: testServiceID,
				Version: tv2.Number,
				Name:    "test-backend",
			})
		})
	}()

	// Diff should mot be empty
	record(t, "diff/get_again", func(c *Client) {
		d, err = c.GetDiff(&GetDiffInput{
			Service: testServiceID,
			From:    tv1.Number,
			To:      tv2.Number,
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
	_, err = testClient.GetDiff(&GetDiffInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDiff(&GetDiffInput{
		Service: "foo",
		From:    0,
	})
	if err != ErrMissingFrom {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDiff(&GetDiffInput{
		Service: "foo",
		From:    1,
		To:      0,
	})
	if err != ErrMissingTo {
		t.Errorf("bad error: %s", err)
	}
}
