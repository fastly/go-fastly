package fastly

import (
	"testing"
	"time"
)

func TestClient_GetDomainMetricsForService(t *testing.T) {
	t.Parallel()

	now := time.Now()
	year, month, day := now.Date()
	end := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, 0, -1)
	var err error

	Record(t, "domain_inspector/metrics_for_service", func(c *Client) {
		_, err = c.GetDomainMetricsForService(&GetDomainMetricsInput{
			ServiceID:   TestDeliveryServiceID,
			Start:       &start,
			End:         &end,
			Domains:     []string{"domain_1.com", "domain_2.com"},
			Datacenters: []string{"SJC", "STP"},
			Metrics:     []string{"resp_body_bytes", "status_2xx"},
			GroupBy:     []string{"domain"},
			Downsample:  ToPointer("hour"),
			Regions:     []string{"usa"},
			Limit:       ToPointer(10),
			Cursor:      ToPointer(""),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
