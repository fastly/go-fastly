package paths

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths/rules"
)

func TestClient_Path(t *testing.T) {
	t.Parallel()

	var err error

	// Setup: a path always belongs to a routing config.
	var rc *routingconfigs.Data
	fastly.Record(t, "setup_config", func(c *fastly.Client) {
		rc, err = routingconfigs.Create(context.TODO(), c, &routingconfigs.CreateInput{
			Name: new("gofastly-sdk-testing-paths-config"),
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

	// Create
	var d *Data
	fastly.Record(t, "create", func(c *fastly.Client) {
		d, err = Create(context.TODO(), c, &CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			Path:            new("/fiesta"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Path != "/fiesta" {
		t.Errorf("bad path: %v", d.Path)
	}
	if d.PathID == "" {
		t.Errorf("bad path id: %v", d.PathID)
	}

	// List
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
		t.Errorf("bad paths list: %v", cl)
	}

	// Get
	var gd *Data
	fastly.Record(t, "get", func(c *fastly.Client) {
		gd, err = Get(context.TODO(), c, &GetInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &d.PathID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gd.Path != "/fiesta" {
		t.Errorf("bad path: %q", gd.Path)
	}

	// Update
	var ud *Data
	fastly.Record(t, "update", func(c *fastly.Client) {
		ud, err = Update(context.TODO(), c, &UpdateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &d.PathID,
			Path:            new("/fiesta2"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ud.Path != "/fiesta2" {
		t.Errorf("bad path: %q", ud.Path)
	}

	// Delete
	fastly.Record(t, "delete", func(c *fastly.Client) {
		err = Delete(context.TODO(), c, &DeleteInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &d.PathID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetPath_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RoutingConfigID: new("abc"),
		PathID:          nil,
	})
	if !errors.Is(err, fastly.ErrMissingPathID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RoutingConfigID: nil,
		PathID:          new("abc"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RoutingConfigID: new(""),
		PathID:          new("abc"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RoutingConfigID: new("abc"),
		PathID:          new(""),
	})
	if !errors.Is(err, fastly.ErrMissingPathID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreatePath_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		RoutingConfigID: nil,
		Path:            new("/"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		RoutingConfigID: new("abc"),
		Path:            nil,
	})
	if !errors.Is(err, fastly.ErrMissingPath) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdatePath_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		RoutingConfigID: nil,
		PathID:          new("abc"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		RoutingConfigID: new("abc"),
		PathID:          nil,
	})
	if !errors.Is(err, fastly.ErrMissingPathID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeletePath_validation(t *testing.T) {
	err := Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		RoutingConfigID: nil,
		PathID:          new("abc"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		RoutingConfigID: new("abc"),
		PathID:          nil,
	})
	if !errors.Is(err, fastly.ErrMissingPathID) {
		t.Errorf("bad error: %s", err)
	}
}

// TestClient_Path_ListReflectsActiveVersion: the docs
// for this endpoint state it "returns paths from the active version if one
// exists, otherwise from the draft," but in practice, once List has been
// called on a routing config with no version yet (a legitimate 404), that
// 404 response is cached and continues being served for the same URL even
// after paths are created and the routing config is activated - despite the
// cached response itself declaring Cache-Control: no-store. This test
// reproduces that trigger (a List call while the config is still empty) and
// then asserts the documented behavior; it is expected to FAIL until the
// caching bug is fixed.
func TestClient_Path_ListReflectsActiveVersion(t *testing.T) {
	t.Parallel()

	var err error

	var rc *routingconfigs.Data
	fastly.Record(t, "active_list_bug_setup_config", func(c *fastly.Client) {
		rc, err = routingconfigs.Create(context.TODO(), c, &routingconfigs.CreateInput{
			Name: new("gofastly-sdk-testing-active-list-bug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		fastly.Record(t, "active_list_bug_teardown_config", func(c *fastly.Client) {
			_ = routingconfigs.Delete(context.TODO(), c, &routingconfigs.DeleteInput{
				RoutingConfigID: &rc.RoutingConfigID,
				Force:           new(true),
			})
		})
	}()

	// Trigger: list paths while the routing config still has none. This is a
	// legitimate 404 ("No version found for config"), but it's what ends up
	// cached for this URL.
	fastly.Record(t, "active_list_bug_list_before", func(c *fastly.Client) {
		_, err = List(context.TODO(), c, &ListInput{RoutingConfigID: &rc.RoutingConfigID})
	})
	if err == nil {
		t.Fatal("expected an error listing paths on a routing config with no version yet")
	}

	var p *Data
	fastly.Record(t, "active_list_bug_create_path", func(c *fastly.Client) {
		p, err = Create(context.TODO(), c, &CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			Path:            new("/"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "active_list_bug_create_rule", func(c *fastly.Client) {
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

	fastly.Record(t, "active_list_bug_activate", func(c *fastly.Client) {
		_, err = routingconfigs.Activate(context.TODO(), c, &routingconfigs.ActivateInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Per the documented behavior, this should now return the path we just
	// activated. As of CDTOOL-1743 it instead returns the cached 404 from the
	// "list before" call above.
	var cl []Data
	fastly.Record(t, "active_list_bug_list_after", func(c *fastly.Client) {
		cl, err = List(context.TODO(), c, &ListInput{RoutingConfigID: &rc.RoutingConfigID})
	})
	if err != nil {
		t.Fatalf("List after activation should succeed per documented behavior, got error: %v", err)
	}
	if len(cl) != 1 {
		t.Errorf("expected 1 active path after activation, got %d: %v", len(cl), cl)
	}
}
