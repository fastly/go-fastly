package fastly

import (
	"context"
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
		d, err = c.CreateDomain(context.TODO(), &CreateDomainInput{
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
		d2, err = c.CreateDomain(context.TODO(), &CreateDomainInput{
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
			_ = c.DeleteDomain(context.TODO(), &DeleteDomainInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           domain1,
			})

			_ = c.DeleteDomain(context.TODO(), &DeleteDomainInput{
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
		ds, err = c.ListDomains(context.TODO(), &ListDomainsInput{
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

	// List with include=staging_ips
	var dsi []*Domain
	Record(t, "domains/list_with_staging_ips", func(c *Client) {
		dsi, err = c.ListDomains(context.TODO(), &ListDomainsInput{
			ServiceID:         TestDeliveryServiceID,
			ServiceVersion:    *tv.Number,
			IncludeStagingIPs: true,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(dsi) < 2 {
		t.Errorf("bad domains: %v", dsi)
	}
	for _, d := range dsi {
		if d.StagingIP == nil {
			t.Errorf("expected staging_ip to be populated for domain %q", *d.Name)
		}
	}

	// Get
	var nd *Domain
	Record(t, "domains/get", func(c *Client) {
		nd, err = c.GetDomain(context.TODO(), &GetDomainInput{
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
		ud, err = c.UpdateDomain(context.TODO(), &UpdateDomainInput{
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
		vd, err = c.ValidateDomain(context.TODO(), &ValidateDomainInput{
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
		vds, err = c.ValidateAllDomains(context.TODO(), &ValidateAllDomainsInput{
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
		err = c.DeleteDomain(context.TODO(), &DeleteDomainInput{
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
	_, err = TestClient.ListDomains(context.TODO(), &ListDomainsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListDomains(context.TODO(), &ListDomainsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDomain_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateDomain(context.TODO(), &CreateDomainInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateDomain(context.TODO(), &CreateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDomain_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetDomain(context.TODO(), &GetDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDomain(context.TODO(), &GetDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetDomain(context.TODO(), &GetDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDomain_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateDomain(context.TODO(), &UpdateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDomain(context.TODO(), &UpdateDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateDomain(context.TODO(), &UpdateDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDomain_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteDomain(context.TODO(), &DeleteDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDomain(context.TODO(), &DeleteDomainInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteDomain(context.TODO(), &DeleteDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ValidateDomain_validation(t *testing.T) {
	var err error

	_, err = TestClient.ValidateDomain(context.TODO(), &ValidateDomainInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ValidateDomain(context.TODO(), &ValidateDomainInput{
		Name:           "test",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ValidateDomain(context.TODO(), &ValidateDomainInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
