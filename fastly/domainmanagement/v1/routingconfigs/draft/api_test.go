package draft

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths/rules"
)

func TestClient_Draft(t *testing.T) {
	t.Parallel()

	var err error

	// Setup: create a routing config with a single path/default-rule, then
	// activate it so subsequent edits produce a draft that differs from the
	// active version.
	var rc *routingconfigs.Data
	fastly.Record(t, "setup_config", func(c *fastly.Client) {
		rc, err = routingconfigs.Create(context.TODO(), c, &routingconfigs.CreateInput{
			Name: new("gofastly-sdk-testing-draft-config"),
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

	var p1 *paths.Data
	fastly.Record(t, "setup_path_1", func(c *fastly.Client) {
		p1, err = paths.Create(context.TODO(), c, &paths.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			Path:            new("/"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "setup_rule_1", func(c *fastly.Client) {
		_, err = rules.Create(context.TODO(), c, &rules.CreateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			PathID:          &p1.PathID,
			Action: &rules.Action{
				Type:  "service",
				Value: fastly.TestDeliveryServiceID,
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	fastly.Record(t, "setup_activate", func(c *fastly.Client) {
		_, err = routingconfigs.Activate(context.TODO(), c, &routingconfigs.ActivateInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Add a second path (with its own default rule) to create a draft that
	// differs from the active version.
	var p2 *paths.Data
	fastly.Record(t, "setup_path_2", func(c *fastly.Client) {
		p2, err = paths.Create(context.TODO(), c, &paths.CreateInput{
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
			PathID:          &p2.PathID,
			Action: &rules.Action{
				Type:  "service",
				Value: fastly.TestDeliveryServiceID,
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// GetDiff
	var diff *Diff
	fastly.Record(t, "get_diff", func(c *fastly.Client) {
		diff, err = GetDiff(context.TODO(), c, &GetDiffInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(diff.Added) != 1 {
		t.Errorf("bad diff added: %v", diff.Added)
	}

	// Update (set a comment on the draft)
	var d *Data
	fastly.Record(t, "update", func(c *fastly.Client) {
		d, err = Update(context.TODO(), c, &UpdateInput{
			RoutingConfigID: &rc.RoutingConfigID,
			Comment:         new("testing comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.Comment != "testing comment" {
		t.Errorf("bad comment: %q", d.Comment)
	}

	// Delete (discard the draft)
	fastly.Record(t, "delete", func(c *fastly.Client) {
		err = Delete(context.TODO(), c, &DeleteInput{
			RoutingConfigID: &rc.RoutingConfigID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetDiff_validation(t *testing.T) {
	_, err := GetDiff(context.TODO(), fastly.TestClient, &GetDiffInput{
		RoutingConfigID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateDraft_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		RoutingConfigID: nil,
		Comment:         new("x"),
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		RoutingConfigID: new("abc"),
		Comment:         nil,
	})
	if !errors.Is(err, fastly.ErrMissingComment) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteDraft_validation(t *testing.T) {
	err := Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		RoutingConfigID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRoutingConfigID) {
		t.Errorf("bad error: %s", err)
	}
}
