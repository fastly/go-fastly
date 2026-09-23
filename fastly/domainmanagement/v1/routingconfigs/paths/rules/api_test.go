package rules

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths"
)

func TestClient_Rule(t *testing.T) {
	t.Parallel()

	var err error

	// Setup: a rule always belongs to a path, which belongs to a routing
	// config.
	var rc *routingconfigs.Data
	fastly.Record(t, "setup_config", func(c *fastly.Client) {
		rc, err = routingconfigs.Create(context.TODO(), c, &routingconfigs.CreateInput{
			Name: new("gofastly-sdk-testing-rules-config"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		fastly.Record(t, "teardown_config", func(c *fastly.Client) {
			_ = routingconfigs.Delete(context.TODO(), c, &routingconfigs.DeleteInput{
				RoutingConfigID: &rc.RoutingConfigID,
			})
		})
	}()

	var p *paths.Data
	fastly.Record(t, "setup_path", func(c *fastly.Client) {
		p, err = paths.Create(context.TODO(), c, &paths.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			Path:            new("/fiesta"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a conditional rule.
	var d *Data
	fastly.Record(t, "create", func(c *fastly.Client) {
		d, err = Create(context.TODO(), c, &CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &p.PathID,
			Action: &Action{
				Type:  "service",
				Value: fastly.TestDeliveryServiceID,
			},
			Conditions: []Condition{
				{
					Type:     "header",
					Key:      new("X-Test"),
					Operator: "equals",
					Value:    "1",
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.IsDefault {
		t.Errorf("expected non-default rule")
	}
	if len(d.Conditions) != 1 {
		t.Errorf("bad conditions: %v", d.Conditions)
	}

	// Create the default (catch-all) rule for the path.
	var defaultRule *Data
	fastly.Record(t, "create_default", func(c *fastly.Client) {
		defaultRule, err = Create(context.TODO(), c, &CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &p.PathID,
			Action: &Action{
				Type:  "service",
				Value: fastly.TestDeliveryServiceID,
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !defaultRule.IsDefault {
		t.Errorf("expected default rule")
	}

	// List
	var cl *Collection
	fastly.Record(t, "list", func(c *fastly.Client) {
		cl, err = List(context.TODO(), c, &ListInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &p.PathID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Data) != 2 {
		t.Errorf("bad rules list: %v", cl)
	}

	// Get
	var gd *Data
	fastly.Record(t, "get", func(c *fastly.Client) {
		gd, err = Get(context.TODO(), c, &GetInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &p.PathID,
			RuleID:          &d.RuleID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gd.RuleID != d.RuleID {
		t.Errorf("bad rule id: %q", gd.RuleID)
	}

	// Update
	var ud *Data
	fastly.Record(t, "update", func(c *fastly.Client) {
		ud, err = Update(context.TODO(), c, &UpdateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &p.PathID,
			RuleID:          &d.RuleID,
			Conditions: &[]Condition{
				{
					Type:     "header",
					Key:      new("X-Test"),
					Operator: "equals",
					Value:    "2",
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ud.Conditions) != 1 || ud.Conditions[0].Value != "2" {
		t.Errorf("bad conditions: %v", ud.Conditions)
	}

	// Delete
	fastly.Record(t, "delete", func(c *fastly.Client) {
		err = Delete(context.TODO(), c, &DeleteInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &p.PathID,
			RuleID:          &d.RuleID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetRule_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RoutingConfigID: nil,
		PathID:          new("abc"),
		RuleID:          new("abc"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RoutingConfigID: new("abc"),
		PathID:          nil,
		RuleID:          new("abc"),
	})
	if !errors.Is(err, fastly.ErrMissingPathID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RoutingConfigID: new("abc"),
		PathID:          new("abc"),
		RuleID:          nil,
	})
	if !errors.Is(err, fastly.ErrMissingRuleID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateRule_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		RoutingConfigID: nil,
		PathID:          new("abc"),
		Action:          &Action{Type: "service", Value: "abc"},
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		RoutingConfigID: new("abc"),
		PathID:          nil,
		Action:          &Action{Type: "service", Value: "abc"},
	})
	if !errors.Is(err, fastly.ErrMissingPathID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		RoutingConfigID: new("abc"),
		PathID:          new("abc"),
		Action:          nil,
	})
	if !errors.Is(err, fastly.ErrMissingAction) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteRule_validation(t *testing.T) {
	err := Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		RoutingConfigID: nil,
		PathID:          new("abc"),
		RuleID:          new("abc"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}
