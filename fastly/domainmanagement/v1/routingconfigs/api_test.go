package routingconfigs

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths/rules"
)

func TestClient_RoutingConfig(t *testing.T) {
	t.Parallel()

	var err error
	name := "gofastly-sdk-testing-routing-config"

	// Create
	var d *Data
	fastly.Record(t, "create", func(c *fastly.Client) {
		d, err = Create(context.TODO(), c, &CreateInput{
			Name: new(name),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != name {
		t.Errorf("bad name: %v", d.Name)
	}
	if d.RoutingConfigID == "" {
		t.Errorf("bad routing config id: %v", d.RoutingConfigID)
	}

	// A routing config cannot be activated until it has at least one path
	// with a default (catch-all) rule, so set that up before activating.
	var p *paths.Data
	fastly.Record(t, "create_path", func(c *fastly.Client) {
		p, err = paths.Create(context.TODO(), c, &paths.CreateInput{
			RoutingConfigID: &d.RoutingConfigID,
			Path:            new("/"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "create_rule", func(c *fastly.Client) {
		_, err = rules.Create(context.TODO(), c, &rules.CreateInput{
			RoutingConfigID: &d.RoutingConfigID,
			PathID:          &p.PathID,
			Action: &rules.Action{
				Type:  "service",
				Value: fastly.TestDeliveryServiceID,
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// List
	var cl *Collection
	fastly.Record(t, "list", func(c *fastly.Client) {
		cl, err = List(context.TODO(), c, &ListInput{
			Limit: new(10),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Data) < 1 {
		t.Errorf("bad routing configs list: %v", cl)
	}

	// Get
	var gd *Data
	fastly.Record(t, "get", func(c *fastly.Client) {
		gd, err = Get(context.TODO(), c, &GetInput{
			RoutingConfigID: &d.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gd.Name != name {
		t.Errorf("bad name: %q", gd.Name)
	}

	// Activate
	var ad *Data
	fastly.Record(t, "activate", func(c *fastly.Client) {
		ad, err = Activate(context.TODO(), c, &ActivateInput{
			RoutingConfigID: &d.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ad.State != "active" {
		t.Errorf("bad state: %q", ad.State)
	}

	// Deactivate
	var dd *Data
	fastly.Record(t, "deactivate", func(c *fastly.Client) {
		dd, err = Deactivate(context.TODO(), c, &DeactivateInput{
			RoutingConfigID: &d.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if dd.State != "draft-only" {
		t.Errorf("bad state: %q", dd.State)
	}

	// Delete
	fastly.Record(t, "delete", func(c *fastly.Client) {
		err = Delete(context.TODO(), c, &DeleteInput{
			RoutingConfigID: &d.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetRoutingConfig_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RoutingConfigID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateRoutingConfig_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Name: nil,
	})
	if !errors.Is(err, fastly.ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteRoutingConfig_validation(t *testing.T) {
	err := Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		RoutingConfigID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ActivateRoutingConfig_validation(t *testing.T) {
	_, err := Activate(context.TODO(), fastly.TestClient, &ActivateInput{
		RoutingConfigID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeactivateRoutingConfig_validation(t *testing.T) {
	_, err := Deactivate(context.TODO(), fastly.TestClient, &DeactivateInput{
		RoutingConfigID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}
