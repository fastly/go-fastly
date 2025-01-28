package fastly

import (
	"errors"
	"testing"
)

func TestClient_GetRealtimeStats_validation(t *testing.T) {
	var err error
	_, err = TestStatsClient.GetRealtimeStats(&GetRealtimeStatsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestStatsClient_GetRealtimeStats(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	RecordRealtimeStats(t, "realtime_stats/get", func(c *RTSClient) {
		_, err = c.GetRealtimeStats(&GetRealtimeStatsInput{
			ServiceID: TestDeliveryServiceID,
			Timestamp: 0,
			Limit:     ToPointer(uint32(3)),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatsClient_GetRealtimeStatsJSON(t *testing.T) {
	t.Parallel()

	var ret struct {
		RenameTimestamp uint64 `json:"Timestamp"`
	}

	var err error
	RecordRealtimeStats(t, "realtime_stats/get", func(c *RTSClient) {
		err = c.GetRealtimeStatsJSON(&GetRealtimeStatsInput{
			ServiceID: TestDeliveryServiceID,
			Timestamp: 0,
			Limit:     ToPointer(uint32(3)),
		}, &ret)
	})
	if err != nil {
		t.Fatal(err)
	}

	if ret.RenameTimestamp == 0 {
		t.Fatalf("got RenameTimestamp=%d, want nonzero", ret.RenameTimestamp)
	}
}
