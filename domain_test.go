package fastly

import "testing"

func TestClient_Domains(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "domains/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var d *Domain
	record(t, "domains/create", func(c *Client) {
		d, err = c.CreateDomain(&CreateDomainInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "integ-test.go-fastly.com",
			Comment: "comment",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "domains/cleanup", func(c *Client) {
			c.DeleteDomain(&DeleteDomainInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "integ-test.go-fastly.com",
			})

			c.DeleteDomain(&DeleteDomainInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-integ-test.go-fastly.com",
			})
		})
	}()

	if d.Name != "integ-test.go-fastly.com" {
		t.Errorf("bad name: %q", d.Name)
	}
	if d.Comment != "comment" {
		t.Errorf("bad comment: %q", d.Comment)
	}

	// List
	var ds []*Domain
	record(t, "domains/list", func(c *Client) {
		ds, err = c.ListDomains(&ListDomainsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) < 1 {
		t.Errorf("bad domains: %v", ds)
	}

	// Get
	var nd *Domain
	record(t, "domains/get", func(c *Client) {
		nd, err = c.GetDomain(&GetDomainInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "integ-test.go-fastly.com",
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
			Service: testServiceID,
			Version: tv.Number,
			Name:    "integ-test.go-fastly.com",
			NewName: "new-integ-test.go-fastly.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.Name != "new-integ-test.go-fastly.com" {
		t.Errorf("bad name: %q", ud.Name)
	}

	// Delete
	record(t, "domains/delete", func(c *Client) {
		err = c.DeleteDomain(&DeleteDomainInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-integ-test.go-fastly.com",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListDomains_validation(t *testing.T) {
	var err error
	_, err = testClient.ListDomains(&ListDomainsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListDomains(&ListDomainsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateDomain_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateDomain(&CreateDomainInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateDomain(&CreateDomainInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetDomain_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDomain(&GetDomainInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDomain(&GetDomainInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDomain(&GetDomainInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDomain_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateDomain(&UpdateDomainInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDomain(&UpdateDomainInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateDomain(&UpdateDomainInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDomain_validation(t *testing.T) {
	var err error
	err = testClient.DeleteDomain(&DeleteDomainInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDomain(&DeleteDomainInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteDomain(&DeleteDomainInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
