package fastly

import (
	"testing"
	"time"
)

func TestClient_GetOriginMetricsForService(t *testing.T) {
	t.Parallel()

	now := time.Now()
	year, month, day := now.Date()
	end := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, 0, -1)
	var err error

	Record(t, "origin_inspector/metrics_for_service", func(c *Client) {
		_, err = c.GetOriginMetricsForService(&GetOriginMetricsInput{
			Cursor:      ToPointer(""),
			Datacenters: []string{"LHR", "JFK"},
			Downsample:  ToPointer("day"),
			End:         &end,
			GroupBy:     []string{"host"},
			Hosts:       []string{"host01"},
			Metrics:     []string{"responses", "status_2xx"},
			Regions:     []string{"europe", "usa"},
			ServiceID:   TestDeliveryServiceID,
			Start:       &start,
			Limit:       ToPointer(150),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
