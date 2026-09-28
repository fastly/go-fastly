package paths

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs"
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
