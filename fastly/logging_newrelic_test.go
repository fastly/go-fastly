package fastly

import (
	"errors"
	"testing"
)

func TestClient_NewRelic(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "newrelic/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var newRelicResp1, newRelicResp2 *NewRelic
	Record(t, "newrelic/create", func(c *Client) {
		newRelicResp1, err = c.CreateNewRelic(&CreateNewRelicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-newrelic"),
			Token:          ToPointer("abcd1234"),
			Format:         ToPointer("format"),
			Placement:      ToPointer("waf_debug"),
			Region:         ToPointer("us"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "newrelic/create2", func(c *Client) {
		newRelicResp2, err = c.CreateNewRelic(&CreateNewRelicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-newrelic-2"),
			Token:          ToPointer("abcd1234"),
			Format:         ToPointer("format"),
			Placement:      ToPointer("waf_debug"),
			Region:         ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail due to an invalid region
	Record(t, "newrelic/create3", func(c *Client) {
		_, err = c.CreateNewRelic(&CreateNewRelicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-newrelic-3"),
			Token:          ToPointer("abcd1234"),
			Format:         ToPointer("format"),
			Placement:      ToPointer("waf_debug"),
			Region:         ToPointer("abc"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "newrelic/delete", func(c *Client) {
			_ = c.DeleteNewRelic(&DeleteNewRelicInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-newrelic",
			})

			_ = c.DeleteNewRelic(&DeleteNewRelicInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-newrelic-2",
			})

			_ = c.DeleteNewRelic(&DeleteNewRelicInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-newrelic",
			})
		})
	}()

	if *newRelicResp1.Name != "test-newrelic" {
		t.Errorf("bad name: %q", *newRelicResp1.Name)
	}
	if *newRelicResp1.Token != "abcd1234" {
		t.Errorf("bad token: %q", *newRelicResp1.Token)
	}
	if *newRelicResp1.Format != "format" {
		t.Errorf("bad format: %q", *newRelicResp1.Format)
	}
	if *newRelicResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *newRelicResp1.FormatVersion)
	}
	if *newRelicResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *newRelicResp1.Placement)
	}
	if *newRelicResp1.Region != "us" {
		t.Errorf("bad region: %q", *newRelicResp1.Region)
	}
	if *newRelicResp2.Name != "test-newrelic-2" {
		t.Errorf("bad name: %q", *newRelicResp2.Name)
	}
	if *newRelicResp2.Token != "abcd1234" {
		t.Errorf("bad token: %q", *newRelicResp2.Token)
	}
	if *newRelicResp2.Format != "format" {
		t.Errorf("bad format: %q", *newRelicResp2.Format)
	}
	if *newRelicResp2.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *newRelicResp2.FormatVersion)
	}
	if *newRelicResp2.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *newRelicResp2.Placement)
	}
	if *newRelicResp2.Region != "eu" {
		t.Errorf("bad region: %q", *newRelicResp2.Region)
	}

	// List
	var ln []*NewRelic
	Record(t, "newrelic/list", func(c *Client) {
		ln, err = c.ListNewRelic(&ListNewRelicInput{
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
	var newRelicGetResp, newRelicGetResp2 *NewRelic
	Record(t, "newrelic/get", func(c *Client) {
		newRelicGetResp, err = c.GetNewRelic(&GetNewRelicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-newrelic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "newrelic/get2", func(c *Client) {
		newRelicGetResp2, err = c.GetNewRelic(&GetNewRelicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-newrelic-2",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *newRelicResp1.Name != *newRelicGetResp.Name {
		t.Errorf("bad name: %q", *newRelicResp1.Name)
	}
	if *newRelicResp1.Token != *newRelicGetResp.Token {
		t.Errorf("bad token: %q", *newRelicResp1.Token)
	}
	if *newRelicResp1.Format != *newRelicGetResp.Format {
		t.Errorf("bad format: %q", *newRelicResp1.Format)
	}
	if *newRelicResp1.FormatVersion != *newRelicGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", *newRelicResp1.FormatVersion)
	}
	if *newRelicResp1.Placement != *newRelicGetResp.Placement {
		t.Errorf("bad placement: %q", *newRelicResp1.Placement)
	}
	if *newRelicResp1.Region != *newRelicGetResp.Region {
		t.Errorf("bad region: %q", *newRelicResp1.Region)
	}
	if *newRelicResp2.Name != *newRelicGetResp2.Name {
		t.Errorf("bad name: %q", *newRelicResp2.Name)
	}
	if *newRelicResp2.Token != *newRelicGetResp2.Token {
		t.Errorf("bad token: %q", *newRelicResp2.Token)
	}
	if *newRelicResp2.Format != *newRelicGetResp2.Format {
		t.Errorf("bad format: %q", *newRelicResp2.Format)
	}
	if *newRelicResp2.FormatVersion != *newRelicGetResp2.FormatVersion {
		t.Errorf("bad format_version: %q", *newRelicResp2.FormatVersion)
	}
	if *newRelicResp2.Placement != *newRelicGetResp2.Placement {
		t.Errorf("bad placement: %q", *newRelicResp2.Placement)
	}
	if *newRelicResp2.Region != *newRelicGetResp2.Region {
		t.Errorf("bad region: %q", *newRelicResp2.Region)
	}

	// Update
	var newRelicUpdateResp1 *NewRelic
	Record(t, "newrelic/update", func(c *Client) {
		newRelicUpdateResp1, err = c.UpdateNewRelic(&UpdateNewRelicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-newrelic",
			NewName:        ToPointer("new-test-newrelic"),
			FormatVersion:  ToPointer(2),
			Region:         ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail due to an invalid region.
	Record(t, "newrelic/update2", func(c *Client) {
		_, err = c.UpdateNewRelic(&UpdateNewRelicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-newrelic",
			Region:         ToPointer("zz"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	if *newRelicUpdateResp1.Name != "new-test-newrelic" {
		t.Errorf("bad name: %q", *newRelicUpdateResp1.Name)
	}
	if *newRelicUpdateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *newRelicUpdateResp1.FormatVersion)
	}
	if *newRelicUpdateResp1.Region != "eu" {
		t.Errorf("bad region: %q", *newRelicUpdateResp1.Region)
	}

	// Delete
	Record(t, "newrelic/delete", func(c *Client) {
		err = c.DeleteNewRelic(&DeleteNewRelicInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-newrelic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListNewRelic_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListNewRelic(&ListNewRelicInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListNewRelic(&ListNewRelicInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateNewRelic_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateNewRelic(&CreateNewRelicInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateNewRelic(&CreateNewRelicInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetNewRelic_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetNewRelic(&GetNewRelicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetNewRelic(&GetNewRelicInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetNewRelic(&GetNewRelicInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateNewRelic_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateNewRelic(&UpdateNewRelicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateNewRelic(&UpdateNewRelicInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateNewRelic(&UpdateNewRelicInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteNewRelic_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteNewRelic(&DeleteNewRelicInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteNewRelic(&DeleteNewRelicInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteNewRelic(&DeleteNewRelicInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
