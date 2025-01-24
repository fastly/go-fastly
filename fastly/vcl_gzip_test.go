package fastly

import (
	"errors"
	"testing"
)

func TestClient_Gzips(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "gzips/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var gzip *Gzip
	Record(t, "gzips/create", func(c *Client) {
		gzip, err = c.CreateGzip(&CreateGzipInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-gzip"),
			ContentTypes:   ToPointer("text/html text/css"),
			Extensions:     ToPointer("html css"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create omissions (GH-7)
	// NOTE: API should return defaults.
	var gzipomit *Gzip
	Record(t, "gzips/create_omissions", func(c *Client) {
		gzipomit, err = c.CreateGzip(&CreateGzipInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-gzip-omit"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *gzipomit.ContentTypes != "text/html application/x-javascript text/css application/javascript text/javascript application/json application/vnd.ms-fontobject application/x-font-opentype application/x-font-truetype application/x-font-ttf application/xml font/eot font/opentype font/otf image/svg+xml image/vnd.microsoft.icon text/plain text/xml" {
		t.Errorf("bad content_types: %q", *gzipomit.ContentTypes)
	}
	if *gzipomit.Extensions != "css js html eot ico otf ttf json" {
		t.Errorf("bad extensions: %q", *gzipomit.Extensions)
	}

	// Ensure deleted
	defer func() {
		Record(t, "gzips/cleanup", func(c *Client) {
			_ = c.DeleteGzip(&DeleteGzipInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-gzip",
			})

			_ = c.DeleteGzip(&DeleteGzipInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-gzip-omit",
			})

			_ = c.DeleteGzip(&DeleteGzipInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-gzip",
			})
		})
	}()

	if *gzip.Name != "test-gzip" {
		t.Errorf("bad name: %q", *gzip.Name)
	}
	if *gzip.ContentTypes != "text/html text/css" {
		t.Errorf("bad content_types: %q", *gzip.ContentTypes)
	}
	if *gzip.Extensions != "html css" {
		t.Errorf("bad extensions: %q", *gzip.Extensions)
	}

	// List
	var gzips []*Gzip
	Record(t, "gzips/list", func(c *Client) {
		gzips, err = c.ListGzips(&ListGzipsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(gzips) < 1 {
		t.Errorf("bad gzips: %v", gzips)
	}

	// Get
	var ngzip *Gzip
	Record(t, "gzips/get", func(c *Client) {
		ngzip, err = c.GetGzip(&GetGzipInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-gzip",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ngzip.Name != *gzip.Name {
		t.Errorf("bad name: %q", *ngzip.Name)
	}
	if *ngzip.ContentTypes != *gzip.ContentTypes {
		t.Errorf("bad content_types: %q", *ngzip.ContentTypes)
	}
	if *ngzip.Extensions != *gzip.Extensions {
		t.Errorf("bad extensions: %q", *ngzip.Extensions)
	}

	// Update
	var ugzip *Gzip
	Record(t, "gzips/update", func(c *Client) {
		ugzip, err = c.UpdateGzip(&UpdateGzipInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-gzip",
			NewName:        ToPointer("new-test-gzip"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ugzip.Name != "new-test-gzip" {
		t.Errorf("bad name: %q", *ugzip.Name)
	}

	// Delete
	Record(t, "gzips/delete", func(c *Client) {
		err = c.DeleteGzip(&DeleteGzipInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-gzip",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListGzips_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListGzips(&ListGzipsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListGzips(&ListGzipsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateGzip_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateGzip(&CreateGzipInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateGzip(&CreateGzipInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetGzip_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetGzip(&GetGzipInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetGzip(&GetGzipInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetGzip(&GetGzipInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateGzip_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateGzip(&UpdateGzipInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateGzip(&UpdateGzipInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateGzip(&UpdateGzipInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteGzip_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteGzip(&DeleteGzipInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteGzip(&DeleteGzipInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteGzip(&DeleteGzipInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
