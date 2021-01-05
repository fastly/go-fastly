package fastly

import (
	"encoding/json"
	"fmt"
)

// RealtimeStatsResponse is a response from Fastly's real-time analytics endpoint
type RealtimeStatsResponse struct {
	Timestamp      uint64          `mapstructure:"Timestamp"`
	Data           []*RealtimeData `mapstructure:"Data"`
	Error          string          `mapstructure:"Error"`
	AggregateDelay uint32          `mapstructure:"AggregateDelay"`
}

// RealtimeData represents combined stats for all Fastly's POPs and aggregate of them.
// It also includes a timestamp of when the stats were recorded
type RealtimeData struct {
	Datacenter map[string]*Stats `mapstructure:"datacenter"`
	Aggregated *Stats            `mapstructure:"aggregated"`
	Recorded   uint64            `mapstructure:"recorded"`
}

// GetRealtimeStatsInput is an input parameter to GetRealtimeStats function
type GetRealtimeStatsInput struct {
	ServiceID string
	Timestamp uint64
	Limit     uint32
}

// GetRealtimeStats returns realtime stats for a service based on the GetRealtimeStatsInput
// parameter. The realtime stats work in a rolling fashion where first request will return
// a timestamp which should be passed to the next call and so on.
// More details at https://developer.fastly.com/reference/api/metrics-stats/realtime/
func (c *RTSClient) GetRealtimeStats(i *GetRealtimeStatsInput) (*RealtimeStatsResponse, error) {
	var resp interface{}
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
func (c *RTSClient) GetRealtimeStatsJSON(i *GetRealtimeStatsInput, dst interface{}) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := fmt.Sprintf("/v1/channel/%s/ts/%d", i.ServiceID, i.Timestamp)

	if i.Limit != 0 {
		path = fmt.Sprintf("%s/limit/%d", path, i.Limit)
	}

	resp, err := c.client.Get(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dst)
}
