package fastly

import "testing"

func TestClient_NewRelic(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "newrelic/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var n *NewRelic
	record(t, "newrelic/create", func(c *Client) {
		n, err = c.CreateNewRelic(&CreateNewRelicInput{
			Service:   testServiceID,
			Version:   tv.Number,
			Name:      String("test-newrelic"),
			Token:     String("abcd1234"),
			Format:    String("format"),
			Placement: String("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "newrelic/delete", func(c *Client) {
			c.DeleteNewRelic(&DeleteNewRelicInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-newrelic",
			})

			c.DeleteNewRelic(&DeleteNewRelicInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-newrelic",
			})
		})
	}()

	if n.Name != "test-newrelic" {
		t.Errorf("bad name: %q", n.Name)
	}
	if n.Token != "abcd1234" {
		t.Errorf("bad token: %q", n.Token)
	}
	if n.Format != "format" {
		t.Errorf("bad format: %q", n.Format)
	}
	if n.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", n.FormatVersion)
	}
	if n.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", n.Placement)
	}

	// List
	var ln []*NewRelic
	record(t, "newrelic/list", func(c *Client) {
		ln, err = c.ListNewRelic(&ListNewRelicInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ln) < 1 {
		t.Errorf("bad newrelics: %v", ln)
	}

	// Get
	var nn *NewRelic
	record(t, "newrelic/get", func(c *Client) {
		nn, err = c.GetNewRelic(&GetNewRelicInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-newrelic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if n.Name != nn.Name {
		t.Errorf("bad name: %q", n.Name)
	}
	if n.Token != nn.Token {
		t.Errorf("bad token: %q", n.Token)
	}
	if n.Format != nn.Format {
		t.Errorf("bad format: %q", n.Format)
	}
	if n.FormatVersion != nn.FormatVersion {
		t.Errorf("bad format_version: %q", n.FormatVersion)
	}
	if n.Placement != nn.Placement {
		t.Errorf("bad placement: %q", n.Placement)
	}

	// Update
	var un *NewRelic
	record(t, "newrelic/update", func(c *Client) {
		un, err = c.UpdateNewRelic(&UpdateNewRelicInput{
			Service:       testServiceID,
			Version:       tv.Number,
			Name:          "test-newrelic",
			NewName:       String("new-test-newrelic"),
			FormatVersion: Uint(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if un.Name != "new-test-newrelic" {
		t.Errorf("bad name: %q", un.Name)
	}
	if un.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", un.FormatVersion)
	}

	// Delete
	record(t, "newrelic/delete", func(c *Client) {
		err = c.DeleteNewRelic(&DeleteNewRelicInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-newrelic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListNewRelic_validation(t *testing.T) {
	var err error
	_, err = testClient.ListNewRelic(&ListNewRelicInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListNewRelic(&ListNewRelicInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateNewRelic_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateNewRelic(&CreateNewRelicInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateNewRelic(&CreateNewRelicInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetNewRelic_validation(t *testing.T) {
	var err error
	_, err = testClient.GetNewRelic(&GetNewRelicInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetNewRelic(&GetNewRelicInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetNewRelic(&GetNewRelicInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateNewRelic_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateNewRelic(&UpdateNewRelicInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateNewRelic(&UpdateNewRelicInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateNewRelic(&UpdateNewRelicInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteNewRelic_validation(t *testing.T) {
	var err error
	err = testClient.DeleteNewRelic(&DeleteNewRelicInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteNewRelic(&DeleteNewRelicInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteNewRelic(&DeleteNewRelicInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
