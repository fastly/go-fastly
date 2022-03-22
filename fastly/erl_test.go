package fastly

import (
	"testing"
)

func TestClient_ERL(t *testing.T) {
	t.Parallel()

	fixtureBase := "erls/"
	testVersion := createTestVersion(t, fixtureBase+"version", testServiceID)

	// Create
	var (
		e *ERL
		err error
	)
	record(t, fixtureBase+"create", func(c *Client) {
		e, err = c.CreateERL(&CreateERLInput{
			ServiceID:          testServiceID,
			ServiceVersion:     testVersion.Number,
			Name:               "test_erl",
			Action:             ERLActionResponse,
			ClientKey:          []string{"req.http.Fastly-Client-IP"},
			HttpMethods:        []ERLHttpMethods{ERLGet, ERLPost},
			PenaltyBoxDuration: 30,
			Response: &ERLResponseType{
				ERLStatus:      429,
				ERLContentType: "application/json",
				ERLContent:     "Too many requests"},
			RpsLimit:   20,
			WindowSize: 10,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteERL(&DeleteERLInput{
				ServiceID:      testServiceID,
				ServiceVersion: testVersion.Number,
				ERLID:          e.ID,
			})
		})
	}()

	if e.Name != "test_erl" {
		t.Errorf("bad name: %q", e.Name)
	}

	// List
	var es []*ERL
	record(t, fixtureBase+"list", func(c *Client) {
		es, err = c.ListERLs(&ListERLsInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	var want, got = 1, len(es)
	if got < want {
		t.Errorf("want %d, got %d", want, got)
	}
	if es[0].Name != "test_erl" {
		t.Errorf("bad name: %q", es[0].Name)
	}

	// Get
	var ne *ERL
	record(t, fixtureBase+"get", func(c *Client) {
		ne, err = c.GetERL(&GetERLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			ERLID:          saveID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if e.ID != ne.ID {
		t.Errorf("bad ID: %q (%q)", e.ID, ge.ID)
	}

	// Update
	var ua *ERL
	record(t, fixtureBase+"update", func(c *Client) {
		ua, err = c.UpdateERL(&UpdateERLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			ID:             saveID,
			Name:           "test_erl",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ua.Name != "test_erl" {
		t.Errorf("Bad name after update %s", ua.Name)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteERL(&DeleteERLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			ERLID:          ne.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListERLs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListERLs(&ListERLsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("error: %s", err)
	}

	_, err = testClient.ListERLs(&ListERLsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateERL_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateERL(&CreateERLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateERL(&CreateERLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetERL_validation(t *testing.T) {
	var err error
	_, err = testClient.GetERL(&GetERLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetERL(&GetERLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetERL(&GetERLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		ERLID:          "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateERL_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateERL(&UpdateERLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})

	if err != ErrMissingServiceID && err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateERL(&UpdateERLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateERL(&UpdateERLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
	_, err = testClient.UpdateERL(&UpdateERLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "acl",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteERL_validation(t *testing.T) {
	var err error
	err = testClient.DeleteERL(&DeleteERLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteERL(&DeleteERLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteERL(&DeleteERLInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		ERLID:          "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
