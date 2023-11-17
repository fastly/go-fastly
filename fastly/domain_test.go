package fastly

import (
	"testing"
)

func TestClient_Domains(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "domains/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// NOTE: Everytime you regenerate the fixtures you'll need to update the
	// domains as they'll potentially be reported as used depending on the
	// service pre-existing.
	domain1 := "integ-test-20221104.go-fastly-1.com"
	domain2 := "integ-test-20221104.go-fastly-2.com"
	domain3 := "integ-test-20221104.go-fastly-3.com"

	// Create
	var d *Domain
	record(t, "domains/create", func(c *Client) {
		d, err = c.CreateDomain(&CreateDomainInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           ToPointer(domain1),
			Comment:        ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	var d2 *Domain
	record(t, "domains/create2", func(c *Client) {
		d2, err = c.CreateDomain(&CreateDomainInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           ToPointer(domain2),
			Comment:        ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "domains/cleanup", func(c *Client) {
			_ = c.DeleteDomain(&DeleteDomainInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           domain1,
			})

			_ = c.DeleteDomain(&DeleteDomainInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           domain3,
			})
		})
	}()

	if d.Name != domain1 {
		t.Errorf("bad name: %q", d.Name)
	}
	if d.Comment != "comment" {
		t.Errorf("bad comment: %q", d.Comment)
	}
	if d2.Name != domain2 {
		t.Errorf("bad name: %q", d.Name)
	}

	// List
	var ds []*Domain
	record(t, "domains/list", func(c *Client) {
		ds, err = c.ListDomains(&ListDomainsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) < 2 {
		t.Errorf("bad domains: %v", ds)
	}

	// Get
	var nd *Domain
	record(t, "domains/get", func(c *Client) {
		nd, err = c.GetDomain(&GetDomainInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           domain1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != nd.Name {
		t.Errorf("bad name: %q (%q)", d.Name, nd.Name)
	}
	if d.Comment != nd.Comment {
		t.Errorf("bad comment: %q (%q)", d.Comment, nd.Comment)
	}

	// Update
	var ud *Domain
	record(t, "domains/update", func(c *Client) {
		ud, err = c.UpdateDomain(&UpdateDomainInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           domain1,
			NewName:        ToPointer(domain3),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.Name != domain3 {
		t.Errorf("bad name: %q", ud.Name)
	}

	// Validate
	var vd *DomainValidationResult
	record(t, "domains/validation", func(c *Client) {
		vd, err = c.ValidateDomain(&ValidateDomainInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           domain3,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if vd.Valid {
		t.Errorf("valid domain unexpected: %q", vd.Metadata.Name)
	}

	var vds []*DomainValidationResult
	record(t, "domains/validate-all", func(c *Client) {
		vds, err = c.ValidateAllDomains(&ValidateAllDomainsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(vds) < 2 {
		t.Errorf("invalid domains: %v", vds)
	}
	for _, d := range vds {
		if d.Valid {
			t.Errorf("valid domain unexpected: %q", d.Metadata.Name)
		}
	}

	// Delete
	record(t, "domains/delete", func(c *Client) {
		err = c.DeleteDomain(&DeleteDomainInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           domain3,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDomains_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDomains(&ListDomainsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDomains(&ListDomainsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDomain_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDomain(&CreateDomainInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDomain(&CreateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDomain_validation(t *testing.T) {
	var err error

	_, err = testClient.GetDomain(&GetDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDomain(&GetDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDomain(&GetDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDomain_validation(t *testing.T) {
	var err error

	_, err = testClient.UpdateDomain(&UpdateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDomain(&UpdateDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDomain(&UpdateDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDomain_validation(t *testing.T) {
	var err error

	err = testClient.DeleteDomain(&DeleteDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDomain(&DeleteDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDomain(&DeleteDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ValidateDomain_validation(t *testing.T) {
	var err error

	_, err = testClient.ValidateDomain(&ValidateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ValidateDomain(&ValidateDomainInput{
		Name:           "test",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ValidateDomain(&ValidateDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
