package fastly

import (
	"errors"
	"testing"
)

func TestClient_DomainsV1(t *testing.T) {
	t.Parallel()

	var err error
	fqdn := "fastly-sdk-gofastly-testing.com"

	// Create
	var dd *DomainsV1Data
	Record(t, "domains_v1/create_domain", func(c *Client) {
		dd, err = c.CreateDomainsV1(&CreateDomainsV1Input{
			FQDN: ToPointer(fqdn),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	// Ensure deleted
	defer func() {
		Record(t, "domains_v1/cleanup_domains", func(c *Client) {
			err = c.DeleteDomainsV1(&DeleteDomainsV1Input{
				DomainID: &dd.DomainID,
			})
		})
	}()

	if dd.FQDN != fqdn {
		t.Errorf("bad fqdn: %v", dd.FQDN)
	}
	if dd.ServiceID != nil {
		t.Errorf("bad service_id: %v", dd.ServiceID)
	}

	// List Definitions
	var ldr *DomainsV1Response
	Record(t, "domains_v1/list_domains", func(c *Client) {
		ldr, err = c.ListDomainsV1(&ListDomainsV1Input{
			// Cursor: ToPointer(""),
			Limit: ToPointer(10),
			FQDN:  ToPointer(dd.FQDN),
			Sort:  ToPointer("fqdn"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ldr.Data) < 1 {
		t.Errorf("bad domains list: %v", ldr)
	}

	// Get
	var gdd *DomainsV1Data
	Record(t, "domains_v1/get_domain", func(c *Client) {
		gdd, err = c.GetDomainsV1(&GetDomainsV1Input{
			DomainID: &dd.DomainID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if dd.FQDN != gdd.FQDN {
		t.Errorf("bad fqdn: %q (%q)", dd.FQDN, gdd.FQDN)
	}

	// Update
	serviceID := "3QSsca5ilwNVpceZ2HJIAH"
	var udd *DomainsV1Data
	Record(t, "domains_v1/update_domain", func(c *Client) {
		udd, err = c.UpdateDomainsV1(&UpdateDomainsV1Input{
			DomainID:  ToPointer(dd.DomainID),
			ServiceID: ToPointer(serviceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if udd.ServiceID == nil || *udd.ServiceID != serviceID {
		t.Errorf("bad service id: %v", udd.ServiceID)
	}

	// Delete
	Record(t, "domains_v1/delete_domain", func(c *Client) {
		err = c.DeleteDomainsV1(&DeleteDomainsV1Input{
			DomainID: &dd.DomainID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetDomainsV1_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetDomainsV1(&GetDomainsV1Input{
		DomainID: nil,
	})
	if !errors.Is(err, ErrMissingDomainID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDomainsV1_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateDomainsV1(&UpdateDomainsV1Input{
		DomainID: nil,
	})
	if !errors.Is(err, ErrMissingDomainID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDomainsV1_validation(t *testing.T) {
	err := TestClient.DeleteDomainsV1(&DeleteDomainsV1Input{
		DomainID: nil,
	})
	if !errors.Is(err, ErrMissingDomainID) {
		t.Errorf("bad error: %s", err)
	}
}
