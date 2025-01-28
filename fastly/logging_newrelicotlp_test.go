package fastly

import (
	"errors"
	"testing"
)

func TestClient_NewRelicOTLP(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "newrelicotlp/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var n *NewRelicOTLP
	Record(t, "newrelicotlp/create", func(c *Client) {
		n, err = c.CreateNewRelicOTLP(&CreateNewRelicOTLPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-newrelicotlp"),
			Token:          ToPointer("abcd1234"),
			URL:            ToPointer("https://example.nr-data.net"),
			Format:         ToPointer("format"),
			Placement:      ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "newrelicotlp/delete", func(c *Client) {
			_ = c.DeleteNewRelicOTLP(&DeleteNewRelicOTLPInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-newrelicotlp",
			})

			_ = c.DeleteNewRelicOTLP(&DeleteNewRelicOTLPInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-newrelicotlp",
			})
		})
	}()

	if *n.Name != "test-newrelicotlp" {
		t.Errorf("bad name: %q", *n.Name)
	}
	if *n.Token != "abcd1234" {
		t.Errorf("bad token: %q", *n.Token)
	}
	if *n.URL != "https://example.nr-data.net" {
		t.Errorf("bad url: %q", *n.URL)
	}
	if *n.Format != "format" {
		t.Errorf("bad format: %q", *n.Format)
	}
	if *n.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *n.FormatVersion)
	}
	if *n.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *n.Placement)
	}

	// List
	var ln []*NewRelicOTLP
	Record(t, "newrelicotlp/list", func(c *Client) {
		ln, err = c.ListNewRelicOTLP(&ListNewRelicOTLPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ln) < 1 {
		t.Errorf("bad newrelics: %v", ln)
	}

	// Get
	var nn *NewRelicOTLP
	Record(t, "newrelicotlp/get", func(c *Client) {
		nn, err = c.GetNewRelicOTLP(&GetNewRelicOTLPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-newrelicotlp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *n.Name != *nn.Name {
		t.Errorf("bad name: %q", *n.Name)
	}
	if *n.Token != *nn.Token {
		t.Errorf("bad token: %q", *n.Token)
	}
	if *n.URL != *nn.URL {
		t.Errorf("bad url: %q", *n.URL)
	}
	if *n.Format != *nn.Format {
		t.Errorf("bad format: %q", *n.Format)
	}
	if *n.FormatVersion != *nn.FormatVersion {
		t.Errorf("bad format_version: %q", *n.FormatVersion)
	}
	if *n.Placement != *nn.Placement {
		t.Errorf("bad placement: %q", *n.Placement)
	}

	// Update
	var un *NewRelicOTLP
	Record(t, "newrelicotlp/update", func(c *Client) {
		un, err = c.UpdateNewRelicOTLP(&UpdateNewRelicOTLPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-newrelicotlp",
			NewName:        ToPointer("new-test-newrelicotlp"),
			FormatVersion:  ToPointer(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *un.Name != "new-test-newrelicotlp" {
		t.Errorf("bad name: %q", *un.Name)
	}
	if *un.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *un.FormatVersion)
	}

	// Delete
	Record(t, "newrelicotlp/delete", func(c *Client) {
		err = c.DeleteNewRelicOTLP(&DeleteNewRelicOTLPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-newrelicotlp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListNewRelicOTLP_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListNewRelicOTLP(&ListNewRelicOTLPInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListNewRelicOTLP(&ListNewRelicOTLPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateNewRelicOTLP_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateNewRelicOTLP(&CreateNewRelicOTLPInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateNewRelicOTLP(&CreateNewRelicOTLPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetNewRelicOTLP_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetNewRelicOTLP(&GetNewRelicOTLPInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetNewRelicOTLP(&GetNewRelicOTLPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetNewRelicOTLP(&GetNewRelicOTLPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateNewRelicOTLP_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateNewRelicOTLP(&UpdateNewRelicOTLPInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateNewRelicOTLP(&UpdateNewRelicOTLPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateNewRelicOTLP(&UpdateNewRelicOTLPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteNewRelicOTLP_validation(t *testing.T) {
	var err error
	err = TestClient.DeleteNewRelicOTLP(&DeleteNewRelicOTLPInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteNewRelicOTLP(&DeleteNewRelicOTLPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteNewRelicOTLP(&DeleteNewRelicOTLPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}
}
