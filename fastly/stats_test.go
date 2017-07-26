package fastly

import "testing"

func TestClient_GetStats(t *testing.T) {
	t.Parallel()

	var err error
	var r *StatsResponse
	record(t, "stats/service_stats", func(c *Client) {
		r, err = c.GetStats(&GetStatsInput{
			Service: testServiceID,
			From:    "10 days ago",
			To:      "now",
			By:      "minute",
			Region:  "europe",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestClient_GetRegions(t *testing.T) {
	t.Parallel()

	var err error
	var r *RegionsResponse
	record(t, "stats/regions", func(c *Client) {
		r, err = c.GetRegions()
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetRegionsUsage(t *testing.T) {
	t.Parallel()

	var err error
	var r *UsageResponse
	record(t, "stats/regions_usage", func(c *Client) {
		r, err = c.GetUsage(&GetUsageInput{
			From:   "10 days ago",
			To:     "now",
			By:     "minute",
			Region: "usa",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetServicesByRegionsUsage(t *testing.T) {
	t.Parallel()

	var err error
	var r *UsageByServiceResponse
	record(t, "stats/services_usage", func(c *Client) {
		r, err = c.GetUsageByService(&GetUsageInput{
			From:   "10 days ago",
			To:     "now",
			By:     "minute",
			Region: "usa",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
