package versions

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths/rules"
)

func TestClient_Version(t *testing.T) {
	t.Parallel()

	var err error

	// Setup: create a routing config, activate it (v1), then edit and
	// activate again (v2) so v1 becomes an inactive, listable version.
	var rc *routingconfigs.Data
	fastly.Record(t, "setup_config", func(c *fastly.Client) {
		rc, err = routingconfigs.Create(context.TODO(), c, &routingconfigs.CreateInput{
			Name: new("gofastly-sdk-testing-versions-config"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		fastly.Record(t, "teardown_config", func(c *fastly.Client) {
			_ = routingconfigs.Delete(context.TODO(), c, &routingconfigs.DeleteInput{
				RoutingConfigID: &rc.RoutingConfigID,
				Force:           new(true),
			})
		})
	}()

	var p *paths.Data
	fastly.Record(t, "setup_path", func(c *fastly.Client) {
		p, err = paths.Create(context.TODO(), c, &paths.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			Path:            new("/"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "setup_rule", func(c *fastly.Client) {
		_, err = rules.Create(context.TODO(), c, &rules.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
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

	fastly.Record(t, "setup_activate_1", func(c *fastly.Client) {
		_, err = routingconfigs.Activate(context.TODO(), c, &routingconfigs.ActivateInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "setup_path_2", func(c *fastly.Client) {
		p, err = paths.Create(context.TODO(), c, &paths.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			Path:            new("/second"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "setup_rule_2", func(c *fastly.Client) {
		_, err = rules.Create(context.TODO(), c, &rules.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
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

	fastly.Record(t, "setup_activate_2", func(c *fastly.Client) {
		_, err = routingconfigs.Activate(context.TODO(), c, &routingconfigs.ActivateInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// List (should show v1 as an inactive version)
	var cl []Data
	fastly.Record(t, "list", func(c *fastly.Client) {
		cl, err = List(context.TODO(), c, &ListInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cl) != 1 {
		t.Fatalf("bad versions list: %v", cl)
	}
	v1 := cl[0]

	// Activate (reactivate v1)
	var ad *routingconfigs.Data
	fastly.Record(t, "activate", func(c *fastly.Client) {
		ad, err = Activate(context.TODO(), c, &ActivateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			VersionID:       &v1.VersionID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ad.State != "active" {
		t.Errorf("bad state: %q", ad.State)
	}

	// DeleteInactive
	//
	// NOTE: as of this writing, the live API returns 404 for this endpoint
	// even when inactive versions exist to delete (verified against a real
	// account, reproducible independent of prior activation state). This
	// assertion reflects that actual behavior rather than the documented
	// success response; update it if/when the backend is fixed.
	fastly.Record(t, "delete_inactive", func(c *fastly.Client) {
		err = DeleteInactive(context.TODO(), c, &DeleteInactiveInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	var httpErr *fastly.HTTPError
	if !errors.As(err, &httpErr) || !httpErr.IsNotFound() {
		t.Fatalf("expected a 404 HTTPError, got: %v", err)
	}
}

func TestClient_ListVersions_validation(t *testing.T) {
	_, err := List(context.TODO(), fastly.TestClient, &ListInput{
		RoutingConfigID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ActivateVersion_validation(t *testing.T) {
	var err error
	_, err = Activate(context.TODO(), fastly.TestClient, &ActivateInput{
		RoutingConfigID: nil,
		VersionID:       new("abc"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Activate(context.TODO(), fastly.TestClient, &ActivateInput{
		RoutingConfigID: new("abc"),
		VersionID:       nil,
	})
	if !errors.Is(err, fastly.ErrMissingVersionID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Activate(context.TODO(), fastly.TestClient, &ActivateInput{
		RoutingConfigID: new("abc"),
		VersionID:       new(""),
	})
	if !errors.Is(err, fastly.ErrMissingVersionID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteInactiveVersions_validation(t *testing.T) {
	err := DeleteInactive(context.TODO(), fastly.TestClient, &DeleteInactiveInput{
		RoutingConfigID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}
