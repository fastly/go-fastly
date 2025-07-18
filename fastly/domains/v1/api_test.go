package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

func TestClient_Domain(t *testing.T) {
	t.Parallel()

	var err error
	fqdn := "fastly-sdk-gofastly-testing.com"

	// Create
	var d *Data
	desc := "my description"
	fastly.Record(t, "create", func(c *fastly.Client) {
		d, err = Create(context.TODO(), c, &CreateInput{
			Description: fastly.ToPointer(desc),
			FQDN:        fastly.ToPointer(fqdn),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Description != desc {
		t.Errorf("bad description: %v", d.Description)
	}
	if d.FQDN != fqdn {
		t.Errorf("bad fqdn: %v", d.FQDN)
	}
	if d.ServiceID != nil {
		t.Errorf("bad service_id: %v", d.ServiceID)
	}

	fastly.Record(t, "create_duplicate", func(c *fastly.Client) {
		_, err = Create(context.TODO(), c, &CreateInput{
			FQDN: fastly.ToPointer(fqdn),
		})
	})
	if err == nil {
		t.Fatal("expected an error and got nil")
	}
	var httpError *fastly.HTTPError
	if !errors.As(err, &httpError) {
		t.Fatalf("unexpected error type: %T", err)
	} else {
		var okErr bool
		for _, he := range httpError.Errors {
			if he.Detail == "fqdn has already been taken" {
				okErr = true
				break
			}
		}
		if !okErr {
			t.Errorf("bad error: %v", err)
		}
	}

	// List Definitions
	var cl *Collection
	fastly.Record(t, "list", func(c *fastly.Client) {
		cl, err = List(context.TODO(), c, &ListInput{
			Limit: fastly.ToPointer(10),
			FQDN:  fastly.ToPointer(d.FQDN),
			Sort:  fastly.ToPointer("fqdn"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Data) < 1 {
		t.Errorf("bad domains list: %v", cl)
	}

	// Get
	var gd *Data
	fastly.Record(t, "get", func(c *fastly.Client) {
		gd, err = Get(context.TODO(), c, &GetInput{
			DomainID: &d.DomainID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Description != gd.Description {
		t.Errorf("bad description: %q (%q)", d.Description, gd.Description)
	}
	if d.FQDN != gd.FQDN {
		t.Errorf("bad fqdn: %q (%q)", d.FQDN, gd.FQDN)
	}

	// Update
	var ud *Data
	descUpdated := "updated description"
	fastly.Record(t, "update", func(c *fastly.Client) {
		ud, err = Update(context.TODO(), c, &UpdateInput{
			Description: fastly.ToPointer(descUpdated),
			DomainID:    fastly.ToPointer(d.DomainID),
			ServiceID:   fastly.ToPointer(fastly.TestDeliveryServiceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.Description != descUpdated {
		t.Errorf("bad description: %q (%q)", descUpdated, ud.Description)
	}
	if ud.ServiceID == nil || *ud.ServiceID != fastly.TestDeliveryServiceID {
		t.Errorf("bad service id: %v", *ud.ServiceID)
	}

	// Delete
	fastly.Record(t, "delete", func(c *fastly.Client) {
		err = Delete(context.TODO(), c, &DeleteInput{
			DomainID: &d.DomainID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetDomain_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		DomainID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingDomainID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDomain_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		DomainID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingDomainID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDomain_validation(t *testing.T) {
	err := Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		DomainID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingDomainID) {
		t.Errorf("bad error: %s", err)
	}
}
