package domains

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths/rules"
)

func TestClient_Domain(t *testing.T) {
	t.Parallel()

	var err error
	fqdn := "1.fastly-sdk-gofastly-testing.com"

	// Create
	var d *Data
	desc := "my description"
	fastly.Record(t, "create", func(c *fastly.Client) {
		d, err = Create(context.TODO(), c, &CreateInput{
			Description: new(desc),
			FQDN:        new(fqdn),
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
			FQDN: new(fqdn),
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
			if strings.Contains(he.Detail, "already been taken") {
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
			Limit: new(10),
			FQDN:  new(d.FQDN),
			Sort:  new("fqdn"),
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

	// Create and activate a routing config to associate with the domain in
	// the Update step below (association requires an active routing config).
	var rc *routingconfigs.Data
	rcName := "gofastly-sdk-testing-domain-routing-config"
	fastly.Record(t, "create_routing_config", func(c *fastly.Client) {
		rc, err = routingconfigs.Create(context.TODO(), c, &routingconfigs.CreateInput{
			Name: new(rcName),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	var rp *paths.Data
	fastly.Record(t, "create_routing_config_path", func(c *fastly.Client) {
		rp, err = paths.Create(context.TODO(), c, &paths.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			Path:            new("/"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "create_routing_config_rule", func(c *fastly.Client) {
		_, err = rules.Create(context.TODO(), c, &rules.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &rp.PathID,
			Action: &rules.Action{
				Type:  "service",
				Value: fastly.TestDeliveryServiceID,
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "activate_routing_config", func(c *fastly.Client) {
		_, err = routingconfigs.Activate(context.TODO(), c, &routingconfigs.ActivateInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Update
	var ud *Data
	descUpdated := "updated description"
	fastly.Record(t, "update", func(c *fastly.Client) {
		ud, err = Update(context.TODO(), c, &UpdateInput{
			Description:            new(descUpdated),
			DomainID:               new(d.DomainID),
			RoutingConfigurationID: fastly.NewNullable(rc.RoutingConfigID),
			ServiceID:              new(fastly.TestDeliveryServiceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.Description != descUpdated {
		t.Errorf("bad description: %q (%q)", descUpdated, ud.Description)
	}
	if ud.ServiceID == nil || *ud.ServiceID != fastly.TestDeliveryServiceID {
		t.Errorf("bad service id: %v", ud.ServiceID)
	}
	if ud.RoutingConfigurationID == nil || *ud.RoutingConfigurationID != rc.RoutingConfigID {
		t.Errorf("bad routing configuration id: %v", ud.RoutingConfigurationID)
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

	// Clean up the routing config created above.
	fastly.Record(t, "deactivate_routing_config", func(c *fastly.Client) {
		_, err = routingconfigs.Deactivate(context.TODO(), c, &routingconfigs.DeactivateInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "delete_routing_config", func(c *fastly.Client) {
		err = routingconfigs.Delete(context.TODO(), c, &routingconfigs.DeleteInput{
			RoutingConfigID: &rc.RoutingConfigID,
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
