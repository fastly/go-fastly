package fastly

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// DomainInspector represents the response format returned for a request to
// the historical Domain Inspector metrics endpoint.
type DomainInspector struct {
	Data   []*DomainData `mapstructure:"data"`
	Meta   *DomainMeta   `mapstructure:"meta"`
	Status *string       `mapstructure:"status"`
}

// DomainData represents the series of values over time for a single
// dimension combination.
type DomainData struct {
	Dimensions map[string]*string `mapstructure:"dimensions"`
	Values     []*DomainMetrics   `mapstructure:"values"`
}

// DomainMetrics represents the possible metrics that can be returned by a call
// to the Domain Inspector endpoints.
type DomainMetrics struct {
	Bandwidth                  *uint64  `mapstructure:"bandwidth"`
	BereqBodyBytes             *uint64  `mapstructure:"bereq_body_bytes"`
	BereqHeaderBytes           *uint64  `mapstructure:"bereq_header_bytes"`
	EdgeHitRatio               *float64 `mapstructure:"edge_hit_ratio"`
	EdgeHitRequests            *uint64  `mapstructure:"edge_hit_requests"`
	EdgeMissRequests           *uint64  `mapstructure:"edge_miss_requests"`
	EdgeRequests               *uint64  `mapstructure:"edge_requests"`
	EdgeRespBodyBytes          *uint64  `mapstructure:"edge_resp_body_bytes"`
	EdgeRespHeaderBytes        *uint64  `mapstructure:"edge_resp_header_bytes"`
	OriginFetchRespBodyBytes   *uint64  `mapstructure:"origin_fetch_resp_body_bytes"`
	OriginFetchRespHeaderBytes *uint64  `mapstructure:"origin_fetch_resp_header_bytes"`
	OriginFetches              *uint64  `mapstructure:"origin_fetches"`
	OriginOffload              *float64 `mapstructure:"origin_offload"`
	OriginStatus1xx            *uint64  `mapstructure:"origin_status_1xx"`
	OriginStatus200            *uint64  `mapstructure:"origin_status_200"`
	OriginStatus204            *uint64  `mapstructure:"origin_status_204"`
	OriginStatus206            *uint64  `mapstructure:"origin_status_206"`
	OriginStatus2xx            *uint64  `mapstructure:"origin_status_2xx"`
	OriginStatus301            *uint64  `mapstructure:"origin_status_301"`
	OriginStatus302            *uint64  `mapstructure:"origin_status_302"`
	OriginStatus304            *uint64  `mapstructure:"origin_status_304"`
	OriginStatus3xx            *uint64  `mapstructure:"origin_status_3xx"`
	OriginStatus400            *uint64  `mapstructure:"origin_status_400"`
	OriginStatus401            *uint64  `mapstructure:"origin_status_401"`
	OriginStatus403            *uint64  `mapstructure:"origin_status_403"`
	OriginStatus404            *uint64  `mapstructure:"origin_status_404"`
	OriginStatus416            *uint64  `mapstructure:"origin_status_416"`
	OriginStatus429            *uint64  `mapstructure:"origin_status_429"`
	OriginStatus4xx            *uint64  `mapstructure:"origin_status_4xx"`
	OriginStatus500            *uint64  `mapstructure:"origin_status_500"`
	OriginStatus501            *uint64  `mapstructure:"origin_status_501"`
	OriginStatus502            *uint64  `mapstructure:"origin_status_502"`
	OriginStatus503            *uint64  `mapstructure:"origin_status_503"`
	OriginStatus504            *uint64  `mapstructure:"origin_status_504"`
	OriginStatus505            *uint64  `mapstructure:"origin_status_505"`
	OriginStatus5xx            *uint64  `mapstructure:"origin_status_5xx"`
	Requests                   *uint64  `mapstructure:"requests"`
	RespBodyBytes              *uint64  `mapstructure:"resp_body_bytes"`
	RespHeaderBytes            *uint64  `mapstructure:"resp_header_bytes"`
	Status1xx                  *uint64  `mapstructure:"status_1xx"`
	Status200                  *uint64  `mapstructure:"status_200"`
	Status204                  *uint64  `mapstructure:"status_204"`
	Status206                  *uint64  `mapstructure:"status_206"`
	Status2xx                  *uint64  `mapstructure:"status_2xx"`
	Status301                  *uint64  `mapstructure:"status_301"`
	Status302                  *uint64  `mapstructure:"status_302"`
	Status304                  *uint64  `mapstructure:"status_304"`
	Status3xx                  *uint64  `mapstructure:"status_3xx"`
	Status400                  *uint64  `mapstructure:"status_400"`
	Status401                  *uint64  `mapstructure:"status_401"`
	Status403                  *uint64  `mapstructure:"status_403"`
	Status404                  *uint64  `mapstructure:"status_404"`
	Status416                  *uint64  `mapstructure:"status_416"`
	Status429                  *uint64  `mapstructure:"status_429"`
	Status4xx                  *uint64  `mapstructure:"status_4xx"`
	Status500                  *uint64  `mapstructure:"status_500"`
	Status501                  *uint64  `mapstructure:"status_501"`
	Status502                  *uint64  `mapstructure:"status_502"`
	Status503                  *uint64  `mapstructure:"status_503"`
	Status504                  *uint64  `mapstructure:"status_504"`
	Status505                  *uint64  `mapstructure:"status_505"`
	Status5xx                  *uint64  `mapstructure:"status_5xx"`
	Timestamp                  *uint64  `mapstructure:"timestamp"`
}

// DomainMeta is the meta section returned for /metrics/domains/... responses.
type DomainMeta struct {
	Downsample *string           `mapstructure:"downsample"`
	End        *string           `mapstructure:"end"`
	Filters    map[string]string `mapstructure:"filters"`
	GroupBy    *string           `mapstructure:"group_by"`
	Limit      *int              `mapstructure:"limit"`
	Metric     *string           `mapstructure:"metric"`
	NextCursor *string           `mapstructure:"next_cursor"`
	Sort       *string           `mapstructure:"sort"`
	Start      *string           `mapstructure:"start"`
}

// GetDomainMetricsInput is the input to a DomainMetrics request.
type GetDomainMetricsInput struct {
	// Cursor is the value from a previous response to retrieve the next page. To request the first page, this should be empty.
	Cursor *string
	// Datacenters limits query to one or more specific POPs.
	Datacenters []string
	// Domains limit query to one or more specific domains.
	Domains []string
	// Downsample is the duration of sample windows.
	Downsample *string
	// End is a valid ISO-8601-formatted date and time, or UNIX timestamp, indicating the exclusive end of the query time range. If not provided, a default is chosen based on the provided downsample value.
	End *time.Time
	// GroupBy is the dimensions to return in the query.
	GroupBy []string
	// Limit is the limit of returned data
	Limit *int
	// Metrics is the metric to retrieve. Up to ten metrics are accepted.
	Metrics []string
	// Regions limits query to one or more specific geographic regions.
	Regions []string
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string
	// Start is a valid ISO-8601-formatted date and time, or UNIX timestamp, indicating the inclusive start of the query time range. If not provided, a default is chosen based on the provided downsample value.
	Start *time.Time
}

// GetDomainMetricsForService retrieves the specified resource.
func (c *Client) GetDomainMetricsForService(i *GetDomainMetricsInput) (*DomainInspector, error) {
	var resp any
	if err := c.GetDomainMetricsForServiceJSON(i, &resp); err != nil {
		return nil, err
	}

	var di *DomainInspector
	if err := decodeMap(resp, &di); err != nil {
		return nil, err
	}
	return di, nil
}

// GetDomainMetricsForServiceJSON retrieves the specified resource.
func (c *Client) GetDomainMetricsForServiceJSON(i *GetDomainMetricsInput, dst any) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := ToSafeURL("metrics", "domains", "services", i.ServiceID)

	ro := &RequestOptions{
		Params: map[string]string{
			"group_by":   strings.Join(i.GroupBy, ","),
			"metric":     strings.Join(i.Metrics, ","),
			"domain":     strings.Join(i.Domains, ","),
			"datacenter": strings.Join(i.Datacenters, ","),
			"region":     strings.Join(i.Regions, ","),
		},
	}
	if i.Cursor != nil {
		ro.Params["cursor"] = *i.Cursor
	}
	if i.Downsample != nil {
		ro.Params["downsample"] = *i.Downsample
	}
	if i.End != nil {
		ro.Params["end"] = strconv.FormatInt(i.End.Unix(), 10)
	}
	if i.Limit != nil {
		ro.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Start != nil {
		ro.Params["start"] = strconv.FormatInt(i.Start.Unix(), 10)
	}

	resp, err := c.Get(path, ro)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dst)
}
