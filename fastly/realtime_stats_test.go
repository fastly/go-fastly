package fastly

import "testing"

func TestClient_GetRealtimeStats_validation(t *testing.T) {
	var err error
	_, err = testStatsClient.GetRealtimeStats(&GetRealtimeStatsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}
}

func TestStatsClient_GetRealtimeStats(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	recordRealtimeStats(t, "realtime_stats/get", func(c *RTSClient) {
		_, err = c.GetRealtimeStats(&GetRealtimeStatsInput{
			Service:   testServiceID,
			Timestamp: 0,
			Limit:     3,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
