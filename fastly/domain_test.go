package fastly

import (
	"errors"
	"testing"
)

func TestClient_Domains(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "domains/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// NOTE: Every time you regenerate the fixtures you'll need to update the
	// domains as they'll potentially be reported as used depending on the
	// service pre-existing.
	domain1 := "integ-test-20221104.go-fastly-1.com"
	domain2 := "integ-test-20221104.go-fastly-2.com"
	domain3 := "integ-test-20221104.go-fastly-3.com"

	// Create
	var d *Domain
	Record(t, "domains/create", func(c *Client) {
		d, err = c.CreateDomain(&CreateDomainInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer(domain1),
			Comment:        ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	var d2 *Domain
	Record(t, "domains/create2", func(c *Client) {
		d2, err = c.CreateDomain(&CreateDomainInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer(domain2),
			Comment:        ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "domains/cleanup", func(c *Client) {
			_ = c.DeleteDomain(&DeleteDomainInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           domain1,
			})

			_ = c.DeleteDomain(&DeleteDomainInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           domain3,
			})
		})
	}()

	if *d.Name != domain1 {
		t.Errorf("bad name: %q", *d.Name)
	}
	if *d.Comment != "comment" {
		t.Errorf("bad comment: %q", *d.Comment)
	}
	if *d2.Name != domain2 {
		t.Errorf("bad name: %q", *d.Name)
	}

	// List
	var ds []*Domain
	Record(t, "domains/list", func(c *Client) {
		ds, err = c.ListDomains(&ListDomainsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "domains/get", func(c *Client) {
		nd, err = c.GetDomain(&GetDomainInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           domain1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *d.Name != *nd.Name {
		t.Errorf("bad name: %q (%q)", *d.Name, *nd.Name)
	}
	if *d.Comment != *nd.Comment {
		t.Errorf("bad comment: %q (%q)", *d.Comment, *nd.Comment)
	}

	// Update
	var ud *Domain
	Record(t, "domains/update", func(c *Client) {
		ud, err = c.UpdateDomain(&UpdateDomainInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           domain1,
			NewName:        ToPointer(domain3),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ud.Name != domain3 {
		t.Errorf("bad name: %q", *ud.Name)
	}

	// Validate
	var vd *DomainValidationResult
	Record(t, "domains/validation", func(c *Client) {
		vd, err = c.ValidateDomain(&ValidateDomainInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           domain3,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *vd.Valid {
		t.Errorf("valid domain unexpected: %q", *vd.Metadata.Name)
	}

	var vds []*DomainValidationResult
	Record(t, "domains/validate-all", func(c *Client) {
		vds, err = c.ValidateAllDomains(&ValidateAllDomainsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(vds) < 2 {
		t.Errorf("invalid domains: %v", vds)
	}
	for _, d := range vds {
		if *d.Valid {
			t.Errorf("valid domain unexpected: %q", *d.Metadata.Name)
		}
	}

	// Delete
	Record(t, "domains/delete", func(c *Client) {
		err = c.DeleteDomain(&DeleteDomainInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           domain3,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDomains_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListDomains(&ListDomainsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListDomains(&ListDomainsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDomain_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateDomain(&CreateDomainInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateDomain(&CreateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDomain_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetDomain(&GetDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDomain(&GetDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDomain(&GetDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDomain_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateDomain(&UpdateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDomain(&UpdateDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDomain(&UpdateDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDomain_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteDomain(&DeleteDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDomain(&DeleteDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDomain(&DeleteDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ValidateDomain_validation(t *testing.T) {
	var err error

	_, err = TestClient.ValidateDomain(&ValidateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ValidateDomain(&ValidateDomainInput{
		Name:           "test",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ValidateDomain(&ValidateDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
