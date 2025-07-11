package fastly

import (
	"context"
	"testing"
)

func TestClient_GetStats(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "stats/service_stats", func(c *Client) {
		_, err = c.GetStats(context.TODO(), &GetStatsInput{
			Service: ToPointer(TestDeliveryServiceID),
			From:    ToPointer("10 days ago"),
			To:      ToPointer("now"),
			By:      ToPointer("minute"),
			Region:  ToPointer("europe"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetStats_ByField(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "stats/service_stats_by_field", func(c *Client) {
		_, err = c.GetStatsField(context.TODO(), &GetStatsInput{
			Field:  ToPointer("bandwidth"),
			From:   ToPointer("1 hour ago"),
			To:     ToPointer("now"),
			By:     ToPointer("minute"),
			Region: ToPointer("europe"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetStats_ByFieldAndService(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "stats/service_stats_by_field_and_service", func(c *Client) {
		_, err = c.GetStats(context.TODO(), &GetStatsInput{
			Service: ToPointer(TestDeliveryServiceID),
			Field:   ToPointer("bandwidth"),
			From:    ToPointer("10 days ago"),
			To:      ToPointer("now"),
			By:      ToPointer("day"),
			Region:  ToPointer("usa"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetStatsJSON(t *testing.T) {
	t.Parallel()

	var ret struct {
		RenameStatus string `json:"status"`
	}

	var err error
	Record(t, "stats/service_stats", func(c *Client) {
		err = c.GetStatsJSON(context.TODO(), &GetStatsInput{
			Service: ToPointer(TestDeliveryServiceID),
			From:    ToPointer("10 days ago"),
			To:      ToPointer("now"),
			By:      ToPointer("minute"),
			Region:  ToPointer("europe"),
		}, &ret)
	})
	if err != nil {
		t.Fatal(err)
	}

	if ret.RenameStatus != "success" {
		t.Fatalf("got RenameStatus=%q, want %q", ret.RenameStatus, "success")
	}
}

func TestClient_GetAggregateJSON(t *testing.T) {
	t.Parallel()

	var ret struct {
		Status string `json:"status"`
	}

	var err error
	Record(t, "stats/aggregate", func(c *Client) {
		err = c.GetAggregateJSON(context.TODO(), &GetAggregateInput{
			From:   ToPointer("15 minutes ago"),
			To:     ToPointer("now"),
			By:     ToPointer("minute"),
			Region: ToPointer("usa"),
		}, &ret)
	})
	if err != nil {
		t.Fatal(err)
	}

	if ret.Status != "success" {
		t.Fatalf("got status=%q, want %q", ret.Status, "success")
	}
}

func TestClient_GetRegions(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "stats/regions", func(c *Client) {
		_, err = c.GetRegions(context.TODO())
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetRegionsUsage(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "stats/regions_usage", func(c *Client) {
		_, err = c.GetUsage(context.TODO(), &GetUsageInput{
			From:   ToPointer("10 days ago"),
			To:     ToPointer("now"),
			By:     ToPointer("minute"),
			Region: ToPointer("usa"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetServicesByRegionsUsage(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "stats/services_usage", func(c *Client) {
		_, err = c.GetUsageByService(context.TODO(), &GetUsageInput{
			From:   ToPointer("10 days ago"),
			To:     ToPointer("now"),
			By:     ToPointer("minute"),
			Region: ToPointer("usa"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
