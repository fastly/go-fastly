package fastly

import (
	"encoding/json"
	"fmt"
)

// RealtimeStatsResponse is a response from Fastly's real-time analytics
// endpoint.
type RealtimeStatsResponse struct {
	// AggregateDelay is how long the system will wait before aggregating messages for each second.
	AggregateDelay *uint32 `mapstructure:"AggregateDelay"`
	// Data is a list of records, each representing one second of time.
	Data  []*RealtimeData `mapstructure:"Data"`
	Error *string         `mapstructure:"Error"`
	// Timestamp is a value to use for subsequent requests.
	Timestamp *uint64 `mapstructure:"Timestamp"`
}

// RealtimeData represents combined stats for all Fastly's POPs and aggregate of them.
// It also includes a timestamp of when the stats were recorded.
type RealtimeData struct {
	// Aggregated aggregates measurements across all Fastly POPs.
	Aggregated *Stats `mapstructure:"aggregated"`
	// Datacenter groups measurements by POP.
	Datacenter map[string]*Stats `mapstructure:"datacenter"`
	// Recorded is the Unix timestamp at which this record's data was generated.
	Recorded *uint64 `mapstructure:"recorded"`
}

// GetRealtimeStatsInput is an input parameter to GetRealtimeStats function.
type GetRealtimeStatsInput struct {
	Limit *uint32
	// ServiceID is the ID of the service (required).
	ServiceID string
	// Timestamp is a value to use for subsequent requests (required).
	Timestamp uint64
}

// GetRealtimeStats returns realtime stats for a service based on the GetRealtimeStatsInput
// parameter. The realtime stats work in a rolling fashion where first request will return
// a timestamp which should be passed to the next call and so on.
// More details at https://developer.fastly.com/reference/api/metrics-stats/realtime/
func (c *RTSClient) GetRealtimeStats(i *GetRealtimeStatsInput) (*RealtimeStatsResponse, error) {
	var resp any
	if err := c.GetRealtimeStatsJSON(i, &resp); err != nil {
		return nil, err
	}

	var s *RealtimeStatsResponse
	if err := decodeMap(resp, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetRealtimeStatsJSON fetches stats and decodes the response directly to the JSON struct dst.
func (c *RTSClient) GetRealtimeStatsJSON(i *GetRealtimeStatsInput, dst any) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := fmt.Sprintf("/v1/channel/%s/ts/%d", i.ServiceID, i.Timestamp)

	if i.Limit != nil {
		path = fmt.Sprintf("%s/limit/%d", path, *i.Limit)
	}

	resp, err := c.client.Get(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dst)
}
