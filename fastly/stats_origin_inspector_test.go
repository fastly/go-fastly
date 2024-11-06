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
			Cursor:      ToPointer(""),
			Datacenters: []string{"LHR", "JFK"},
			Downsample:  ToPointer("day"),
			End:         &end,
			GroupBy:     []string{"host"},
			Hosts:       []string{"host01"},
			Metrics:     []string{"responses", "status_2xx"},
			Regions:     []string{"europe", "usa"},
			ServiceID:   testDeliveryServiceID,
			Start:       &start,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
