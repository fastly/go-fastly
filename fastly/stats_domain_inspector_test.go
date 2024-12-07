package fastly

import (
	"testing"
	"time"
)

func TestClient_GetDomainMetricsForService(t *testing.T) {
	t.Parallel()

	// NOTE: Update this to a recent time when regenerating the test fixtures,
	// otherwise the data may be outside of retention and an error will be
	// returned.
	end := time.Date(2023, 11, 7, 0, 0, 0, 0, time.UTC)
	start := end.Add(-8 * time.Hour)
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
