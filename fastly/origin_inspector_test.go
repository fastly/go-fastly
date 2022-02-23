package fastly

import (
	"testing"
	"time"
)

func TestClient_GetOriginMetricsForService(t *testing.T) {
	t.Parallel()

	// NOTE: Update this to a recent time when regenerating the test fixtures,
	// otherwise the data may be outside of retention and an error will be
	// returned.
	end := time.Date(2022, 2, 14, 0, 0, 0, 0, time.UTC)
	start := end.Add(-2 * 24 * time.Hour)
	var err error
	record(t, "origin_inspector/metrics_for_service", func(c *Client) {
		_, err = c.GetOriginMetricsForService(&GetOriginMetricsInput{
			ServiceID:   testServiceID,
			Start:       start,
			End:         end,
			Hosts:       []string{"host01"},
			Datacenters: []string{"LHR", "JFK"},
			Metrics:     []string{"responses", "status_2xx"},
			GroupBy:     []string{"host"},
			Downsample:  "day",
			Regions:     []string{"europe", "usa"},
			Cursor:      "",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
