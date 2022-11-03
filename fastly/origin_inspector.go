package fastly

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// OriginInspector represents the response format returned for a request to
// the historical Origin Inspector metrics endpoint.
type OriginInspector struct {
	Data   []OriginData `mapstructure:"data"`
	Meta   OriginMeta   `mapstructure:"meta"`
	Status string       `mapstructure:"status"`
}

// OriginData represents the series of values over time for a single
// dimension combination.
type OriginData struct {
	Dimensions map[string]string `mapstructure:"dimensions"`
	Values     []OriginMetrics   `mapstructure:"values"`
}

// OriginMetrics represents the possible metrics that can be returned by a call
// to the Origin Inspector endpoints.
type OriginMetrics struct {
	Latency0to1ms         uint64 `mapstructure:"latency_0_to_1ms"`
	Latency10000to60000ms uint64 `mapstructure:"latency_10000_to_60000ms"`
	Latency1000to5000ms   uint64 `mapstructure:"latency_1000_to_5000ms"`
	Latency100to250ms     uint64 `mapstructure:"latency_100_to_250ms"`
	Latency10to50ms       uint64 `mapstructure:"latency_10_to_50ms"`
	Latency1to5ms         uint64 `mapstructure:"latency_1_to_5ms"`
	Latency250to500ms     uint64 `mapstructure:"latency_250_to_500ms"`
	Latency5000to10000ms  uint64 `mapstructure:"latency_5000_to_10000ms"`
	Latency500to1000ms    uint64 `mapstructure:"latency_500_to_1000ms"`
	Latency50to100ms      uint64 `mapstructure:"latency_50_to_100ms"`
	Latency5to10ms        uint64 `mapstructure:"latency_5_to_10ms"`
	Latency60000ms        uint64 `mapstructure:"latency_60000ms"`
	RespBodyBytes         uint64 `mapstructure:"resp_body_bytes"`
	RespHeaderBytes       uint64 `mapstructure:"resp_header_bytes"`
	Responses             uint64 `mapstructure:"responses"`
	Status1xx             uint64 `mapstructure:"status_1xx"`
	Status200             uint64 `mapstructure:"status_200"`
	Status204             uint64 `mapstructure:"status_204"`
	Status206             uint64 `mapstructure:"status_206"`
	Status2xx             uint64 `mapstructure:"status_2xx"`
	Status301             uint64 `mapstructure:"status_301"`
	Status302             uint64 `mapstructure:"status_302"`
	Status304             uint64 `mapstructure:"status_304"`
	Status3xx             uint64 `mapstructure:"status_3xx"`
	Status400             uint64 `mapstructure:"status_400"`
	Status401             uint64 `mapstructure:"status_401"`
	Status403             uint64 `mapstructure:"status_403"`
	Status404             uint64 `mapstructure:"status_404"`
	Status416             uint64 `mapstructure:"status_416"`
	Status429             uint64 `mapstructure:"status_429"`
	Status4xx             uint64 `mapstructure:"status_4xx"`
	Status500             uint64 `mapstructure:"status_500"`
	Status501             uint64 `mapstructure:"status_501"`
	Status502             uint64 `mapstructure:"status_502"`
	Status503             uint64 `mapstructure:"status_503"`
	Status504             uint64 `mapstructure:"status_504"`
	Status505             uint64 `mapstructure:"status_505"`
	Status5xx             uint64 `mapstructure:"status_5xx"`
	Timestamp             uint64 `mapstructure:"timestamp"`
}

// OriginMeta is the meta section returned for /metrics/origins responses
type OriginMeta struct {
	Downsample string            `mapstructure:"downsample"`
	End        string            `mapstructure:"end"`
	Filters    map[string]string `mapstructure:"filters"`
	GroupBy    string            `mapstructure:"group_by"`
	Limit      int               `mapstructure:"limit"`
	Metric     string            `mapstructure:"metric"`
	NextCursor string            `mapstructure:"next_cursor"`
	Sort       string            `mapstructure:"sort"`
	Start      string            `mapstructure:"start"`
}

// GetOriginMetricsInput is the input to an OriginMetrics request.
type GetOriginMetricsInput struct {
	// Cursor is the value from a previous response to retrieve the next page. To request the first page, this should be empty.
	Cursor string
	// Datacenters limits query to one or more specific POPs.
	Datacenters []string
	// Downsample is the duration of sample windows.
	Downsample string
	// End is a valid ISO-8601-formatted date and time, or UNIX timestamp, indicating the exclusive end of the query time range. If not provided, a default is chosen based on the provided downsample value.
	End time.Time
	// GroupBy is the dimensions to return in the query.
	GroupBy []string
	// Hosts limits query to one or more specific origin hosts.
	Hosts []string
	// Metrics is the metric to retrieve. Up to ten metrics are accepted.
	Metrics []string
	// Regions limits query to one or more specific geographic regions.
	Regions []string
	// ServiceID is an alphanumeric string identifying the service.
	ServiceID string
	// Start is a valid ISO-8601-formatted date and time, or UNIX timestamp, indicating the inclusive start of the query time range. If not provided, a default is chosen based on the provided downsample value.
	Start time.Time
}

// GetOriginMetricsForService retrieves the specified resource.
func (c *Client) GetOriginMetricsForService(i *GetOriginMetricsInput) (*OriginInspector, error) {
	var resp interface{}
	if err := c.GetOriginMetricsForServiceJSON(i, &resp); err != nil {
		return nil, err
	}

	var or *OriginInspector
	if err := decodeMap(resp, &or); err != nil {
		return nil, err
	}
	return or, nil
}

// GetOriginMetricsForServiceJSON retrieves the specified resource.
func (c *Client) GetOriginMetricsForServiceJSON(i *GetOriginMetricsInput, dst interface{}) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	p := "/metrics/origins/services/" + i.ServiceID

	start := ""
	if !i.Start.IsZero() {
		start = strconv.FormatInt(i.Start.Unix(), 10)
	}
	end := ""
	if !i.End.IsZero() {
		end = strconv.FormatInt(i.End.Unix(), 10)
	}

	r, err := c.Get(p, &RequestOptions{
		Params: map[string]string{
			"start":      start,
			"end":        end,
			"downsample": i.Downsample,
			"metric":     strings.Join(i.Metrics, ","),
			"host":       strings.Join(i.Hosts, ","),
			"datacenter": strings.Join(i.Datacenters, ","),
			"region":     strings.Join(i.Regions, ","),
			"cursor":     i.Cursor,
		},
	})
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(dst)
}
