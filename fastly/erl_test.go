package fastly

import (
	"net/http"
	"testing"
)

func TestClient_ERL(t *testing.T) {
	t.Parallel()

	fixtureBase := "erls/"
	testVersion := createTestVersion(t, fixtureBase+"version", testServiceID)

	// Create
	var (
		e   *ERL
		err error
	)
	record(t, fixtureBase+"create", func(c *Client) {
		e, err = c.CreateERL(&CreateERLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           String("test_erl"),
			Action:         ERLActionPtr(ERLActionResponse),
			ClientKey: &[]string{
				"req.http.Fastly-Client-IP",
			},
			HTTPMethods: &[]string{
				http.MethodGet,
				http.MethodPost,
			},
			PenaltyBoxDuration: Int(30),
			Response: &ERLResponseType{
				ERLStatus:      429,
				ERLContentType: "application/json",
				ERLContent:     "Too many requests",
			},
			RpsLimit:   Int(20),
			WindowSize: ERLWindowSizePtr(10),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteERL(&DeleteERLInput{
				ERLID: e.ID,
			})
		})
	}()

	if e.Name != "test_erl" {
		t.Errorf("bad name: %q", e.Name)
	}

	if e.RpsLimit != 20 {
		t.Errorf("wrong value: %q", e.RpsLimit)
	}

	if e.Response.ERLContent != "Too many requests" {
		t.Errorf("want 'Too many requests', got %q", e.Response.ERLContent)
	}

	if e.Response.ERLContentType != "application/json" {
		t.Errorf("want 'application/json', got %q", e.Response.ERLContentType)
	}

	if e.Response.ERLStatus != 429 {
		t.Errorf("want 429, got %q", e.Response.ERLStatus)
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
	want, got := 1, len(es)
	if got < want {
		t.Errorf("want %d, got %d", want, got)
	}
	if es[0].Name != "test_erl" {
		t.Errorf("bad name: %q", es[0].Name)
	}

	// Get
	var ge *ERL
	record(t, fixtureBase+"get", func(c *Client) {
		ge, err = c.GetERL(&GetERLInput{
			ERLID: e.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if e.ID != ge.ID {
		t.Errorf("bad ID: %q (%q)", e.ID, ge.ID)
	}

	// Update
	var ua *ERL
	record(t, fixtureBase+"update", func(c *Client) {
		ua, err = c.UpdateERL(&UpdateERLInput{
			ERLID: e.ID,
			Name:  String("test_erl"),
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
			ERLID: ge.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create logger type
	var elog *ERL
	record(t, fixtureBase+"logger_create", func(c *Client) {
		elog, err = c.CreateERL(&CreateERLInput{
			ServiceID:      testServiceID,
			ServiceVersion: testVersion.Number,
			Name:           String("test_erl"),
			Action:         ERLActionPtr(ERLActionLogOnly),
			// IMPORTANT: API will 400 if LoggerType not set with log_only action.
			LoggerType: ERLLoggerPtr(ERLLogAzureBlob),
			ClientKey: &[]string{
				"req.http.Fastly-Client-IP",
			},
			HTTPMethods: &[]string{
				http.MethodGet,
				http.MethodPost,
			},
			PenaltyBoxDuration: Int(30),
			RpsLimit:           Int(20),
			WindowSize:         ERLWindowSizePtr(10),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		record(t, fixtureBase+"logger_cleanup", func(c *Client) {
			_ = c.DeleteERL(&DeleteERLInput{
				ERLID: elog.ID,
			})
		})
	}()

	if elog.Name != "test_erl" {
		t.Errorf("bad name: %q", elog.Name)
	}

	if elog.LoggerType != ERLLogAzureBlob {
		t.Errorf("bad logger type: %q", elog.LoggerType)
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
		t.Errorf("error: %s", err)
	}
}

func TestClient_CreateERL_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateERL(&CreateERLInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("error: %s", err)
	}

	_, err = testClient.CreateERL(&CreateERLInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("error: %s", err)
	}
}

func TestClient_GetERL_validation(t *testing.T) {
	_, err := testClient.GetERL(&GetERLInput{
		ERLID: "",
	})
	if err != ErrMissingERLID {
		t.Errorf("error: %s", err)
	}
}

func TestClient_UpdateERL_validation(t *testing.T) {
	_, err := testClient.UpdateERL(&UpdateERLInput{
		ERLID: "",
	})
	if err != ErrMissingERLID {
		t.Errorf("error: %s", err)
	}
}

func TestClient_DeleteERL_validation(t *testing.T) {
	err := testClient.DeleteERL(&DeleteERLInput{
		ERLID: "",
	})
	if err != ErrMissingERLID {
		t.Errorf("error: %s", err)
	}
}
