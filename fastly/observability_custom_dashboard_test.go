package fastly

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestClient_ObservabilityCustomDashboards(t *testing.T) {
	t.Parallel()

	cocd := &CreateObservabilityCustomDashboardInput{
		Description: ToPointer("My dashboard is super cool."),
		Name:        "My Cool Dashboard",
		Items: []DashboardItem{{
			DataSource: DashboardDataSource{
				Config: DashboardSourceConfig{
					Metrics: []string{"requests"},
				},
				Type: SourceTypeStatsEdge,
			},
			Span:     4,
			Subtitle: "This is a subtitle",
			Title:    "A Dashboard Item",
			Visualization: DashboardVisualization{
				Config: VisualizationConfig{PlotType: PlotTypeLine},
				Type:   VisualizationTypeChart,
			},
		}},
	}

	var err error

	// Create
	var ocd *ObservabilityCustomDashboard
	Record(t, "observability_custom_dashboards/create_custom_dashboard", func(c *Client) {
		ocd, err = c.CreateObservabilityCustomDashboard(cocd)
	})
	if err != nil {
		t.Fatalf("Error encountered: %v\ninput: %#v", err, cocd)
	}
	// Ensure deleted
	defer func() {
		Record(t, "observability_custom_dashboards/delete_custom_dashboard", func(c *Client) {
			err = c.DeleteObservabilityCustomDashboard(&DeleteObservabilityCustomDashboardInput{
				ID: &ocd.ID,
			})
		})
	}()

	if ocd.Description != "My dashboard is super cool." {
		t.Errorf("bad description. want: %s, got %s", *cocd.Description, ocd.Description)
	}

	if ocd.Name != "My Cool Dashboard" {
		t.Errorf("bad name. want: %s, got %s", cocd.Name, ocd.Name)
	}

	if len(ocd.Items) != 1 {
		t.Errorf("bad items: %v", ocd.Items)
	}

	// List Dashboards
	var ldr *ListDashboardsResponse
	Record(t, "observability_custom_dashboards/list_custom_dashboards", func(c *Client) {
		ldr, err = c.ListObservabilityCustomDashboards(&ListObservabilityCustomDashboardsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ldr.Data) < 1 {
		t.Errorf("bad custom dashboards: %v", ldr)
	}

	// Get
	var gocd *ObservabilityCustomDashboard
	Record(t, "observability_custom_dashboards/get_custom_dashboard", func(c *Client) {
		gocd, err = c.GetObservabilityCustomDashboard(&GetObservabilityCustomDashboardInput{
			ID: &ocd.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gocd.Name != ocd.Name {
		t.Errorf("bad name: %q (%q)", ocd.Name, gocd.Name)
	}

	// Update
	var ucd *ObservabilityCustomDashboard
	items := ocd.Items
	items[0].DataSource.Config.Metrics = []string{"edge_hit_requests"}
	items[0].Visualization.Config.PlotType = PlotTypeSingleMetric
	items[0].Title = "An Updated Dashboard Item"

	items = append(items, DashboardItem{
		DataSource: DashboardDataSource{
			Config: DashboardSourceConfig{
				Metrics: []string{"requests"},
			},
			Type: SourceTypeStatsEdge,
		},
		Span:     4,
		Subtitle: "This is a subtitle",
		Title:    "A New Dashboard Item",
		Visualization: DashboardVisualization{
			Config: VisualizationConfig{PlotType: PlotTypeLine},
			Type:   VisualizationTypeChart,
		},
	})
	Record(t, "observability_custom_dashboards/update_custom_dashboard", func(c *Client) {
		ucd, err = c.UpdateObservabilityCustomDashboard(&UpdateObservabilityCustomDashboardInput{
			Description: ToPointer("My dashboard just got even cooler."),
			ID:          &ocd.ID,
			Items:       &items,
			Name:        ToPointer("My Updated Dashboard"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ucd.Name != "My Updated Dashboard" {
		t.Errorf("bad name: %q (%q)", "My Updated Dashboard", ucd.Name)
	}
	if len(ucd.Items) != 2 {
		t.Errorf("bad items")
	}
	if diff := cmp.Diff(items[0], ucd.Items[0]); diff != "" {
		t.Errorf("dashboard item did not match (-want,+got): %s", diff)
	}
	if diff := cmp.Diff(items[1], ucd.Items[1], cmpopts.IgnoreFields(DashboardItem{}, "ID")); diff != "" {
		t.Errorf("dashboard item did not match (-want,+got): %s", diff)
	}
}
