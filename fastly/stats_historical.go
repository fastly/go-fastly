package fastly

import (
	"encoding/json"
	"fmt"
)

// StatsResponse is a response from the service stats API endpoint.
type StatsResponse struct {
	Data    []*Stats          `mapstructure:"data"`
	Message *string           `mapstructure:"msg"`
	Meta    map[string]string `mapstructure:"meta"`
	Status  *string           `mapstructure:"status"`
}

// Stats represent metrics of a Fastly service.
type Stats struct {
	AttackRequestHeaderBytes  *uint64     `mapstructure:"attack_req_header_bytes"` // Total header bytes received from requests that triggered a WAF rule.
	AttackRequestBodyBytes    *uint64     `mapstructure:"attack_req_body_bytes"`   // Total body bytes received from requests that triggered a WAF rule.
	AttackResponseSynthBytes  *uint64     `mapstructure:"attack_resp_synth_bytes"` // Total bytes delivered for requests that triggered a WAF rule and returned a synthetic response.
	BERequestBodyBytes        *uint64     `mapstructure:"bereq_body_bytes"`        // Total body bytes sent to origin.
	BERequestHeaderbytes      *uint64     `mapstructure:"bereq_header_bytes"`      // Total header bytes sent to origin.
	Bandwidth                 *uint64     `mapstructure:"bandwidth"`               // Total bytes delivered (body_size + header_size).
	BilledBodyBytes           *uint64     `mapstructure:"billed_body_bytes"`
	BilledHeaderBytes         *uint64     `mapstructure:"billed_header_bytes"`
	Errors                    *uint64     `mapstructure:"errors"`                   // Number of cache errors.
	HTTP2                     *uint64     `mapstructure:"http2"`                    // Number of requests received over HTTP2.
	HitRatio                  *float64    `mapstructure:"hit_ratio"`                // Ratio of cache hits to cache misses (between 0 and 1).
	Hits                      *uint64     `mapstructure:"hits"`                     // Number of cache hits.
	HitsTime                  *float64    `mapstructure:"hits_time"`                // Total amount of time spent processing cache hits (in seconds).
	IPv6                      *uint64     `mapstructure:"ipv6"`                     // Number of requests that were received over IPv6.
	ImageOptimizer            *uint64     `mapstructure:"imgopto"`                  // Number of responses that came from the Fastly Image Optimizer service.
	Log                       *uint64     `mapstructure:"log"`                      // Number of log lines sent.
	Miss                      *uint64     `mapstructure:"miss"`                     // Number of cache misses.
	MissHistogram             map[int]int `mapstructure:"miss_histogram"`           // Number of requests to origin in time buckets of 10s of milliseconds
	MissTime                  *float64    `mapstructure:"miss_time"`                // Amount of time spent processing cache misses (in seconds).
	OTFP                      *uint64     `mapstructure:"otfp"`                     // Number of responses that came from the Fastly On-the-Fly Packager for On Demand Streaming service for video-on-demand.
	ObjectSize100k            *uint64     `mapstructure:"object_size_100k"`         // Number of objects served that were between 10KB and 100KB in size.
	ObjectSize100m            *uint64     `mapstructure:"object_size_100m"`         // Number of objects served that were between 10MB and 100MB in size.
	ObjectSize10k             *uint64     `mapstructure:"object_size_10k"`          // Number of objects served that were between 1KB and 10KB in size.
	ObjectSize10m             *uint64     `mapstructure:"object_size_10m"`          // Number of objects served that were between 1MB and 10MB in size.
	ObjectSize1g              *uint64     `mapstructure:"object_size_1g"`           // Number of objects served that were between 100MB and 1GB in size.
	ObjectSize1k              *uint64     `mapstructure:"object_size_1k"`           // Number of objects served that were under 1KB in size.
	ObjectSize1m              *uint64     `mapstructure:"object_size_1m"`           // Number of objects served that were between 100KB and 1MB in size.
	PCI                       *uint64     `mapstructure:"pci"`                      // Number of responses with the PCI flag turned on.
	Pass                      *uint64     `mapstructure:"pass"`                     // Number of requests that passed through the CDN without being cached.
	PassTime                  *float64    `mapstructure:"pass_time"`                // Amount of time spent processing cache passes (in seconds).
	Pipe                      *uint64     `mapstructure:"pipe"`                     // Optional. Pipe operations performed (legacy feature).
	RequestBodyBytes          *uint64     `mapstructure:"req_body_bytes"`           // Total body bytes received.
	RequestHeaderBytes        *uint64     `mapstructure:"req_header_bytes"`         // Total header bytes received.
	Requests                  *uint64     `mapstructure:"requests"`                 // Number of requests processed.
	ResponseBodyBytes         *uint64     `mapstructure:"resp_body_bytes"`          // Total body bytes delivered.
	ResponseHeaderBytes       *uint64     `mapstructure:"resp_header_bytes"`        // Total header bytes delivered.
	Restarts                  *uint64     `mapstructure:"restarts"`                 // Number of restarts performed.
	Shield                    *uint64     `mapstructure:"shield"`                   // Number of requests from shield to origin.
	ShieldResponseBodyBytes   *uint64     `mapstructure:"shield_resp_body_bytes"`   // Total body bytes delivered via a shield.
	ShieldResponseHeaderBytes *uint64     `mapstructure:"shield_resp_header_bytes"` // Total header bytes delivered via a shield.
	Status1xx                 *uint64     `mapstructure:"status_1xx"`               // Number of "Informational" category status codes delivered.
	Status200                 *uint64     `mapstructure:"status_200"`               // Number of responses sent with status code 200 (Success).
	Status204                 *uint64     `mapstructure:"status_204"`               // Number of responses sent with status code 204 (No Content).
	Status206                 *uint64     `mapstructure:"status_206"`               // Number of responses sent with status code 206 (Partial Content).
	Status2xx                 *uint64     `mapstructure:"status_2xx"`               // Number of "Success" status codes delivered.
	Status301                 *uint64     `mapstructure:"status_301"`               // Number of responses sent with status code 301 (Moved Permanently).
	Status302                 *uint64     `mapstructure:"status_302"`               // Number of responses sent with status code 302 (Found).
	Status304                 *uint64     `mapstructure:"status_304"`               // Number of responses sent with status code 304 (Not Modified).
	Status3xx                 *uint64     `mapstructure:"status_3xx"`               // Number of "Redirection" codes delivered.
	Status400                 *uint64     `mapstructure:"status_400"`               // Number of responses sent with status code 400 (Bad Request).
	Status401                 *uint64     `mapstructure:"status_401"`               // Number of responses sent with status code 401 (Unauthorized).
	Status403                 *uint64     `mapstructure:"status_403"`               // Number of responses sent with status code 403 (Forbidden).
	Status404                 *uint64     `mapstructure:"status_404"`               // Number of responses sent with status code 404 (Not Found).
	Status416                 *uint64     `mapstructure:"status_416"`               // Number of responses sent with status code 416 (Range Not Satisfiable).
	Status4xx                 *uint64     `mapstructure:"status_4xx"`               // Number of "Client Error" codes delivered.
	Status500                 *uint64     `mapstructure:"status_500"`               // Number of responses sent with status code 500 (Internal Server Error).
	Status501                 *uint64     `mapstructure:"status_501"`               // Number of responses sent with status code 501 (Not Implemented).
	Status502                 *uint64     `mapstructure:"status_502"`               // Number of responses sent with status code 502 (Bad Gateway).
	Status503                 *uint64     `mapstructure:"status_503"`               // Number of responses sent with status code 503 (Service Unavailable).
	Status504                 *uint64     `mapstructure:"status_504"`               // Number of responses sent with status code 504 (Gateway Timeout).
	Status505                 *uint64     `mapstructure:"status_505"`               // Number of responses sent with status code 505 (HTTP Version Not Supported).
	Status5xx                 *uint64     `mapstructure:"status_5xx"`               // Number of "Server Error" codes delivered.
	Synth                     *uint64     `mapstructure:"synth"`                    // Number of requests that returned synth response.
	TLS                       *uint64     `mapstructure:"tls"`                      // Number of requests that were received over TLS.
	TLSv10                    *uint64     `mapstructure:"tls_v10"`                  // Number of requests received over TLS 1.0.
	TLSv11                    *uint64     `mapstructure:"tls_v11"`                  // Number of requests received over TLS 1.`.
	TLSv12                    *uint64     `mapstructure:"tls_v12"`                  // Number of requests received over TLS 1.2.
	TLSv13                    *uint64     `mapstructure:"tls_v13"`                  // Number of requests received over TLS 1.3.
	Uncachable                *uint64     `mapstructure:"uncachable"`               // Number of requests that were designated uncachable.
	Video                     *uint64     `mapstructure:"video"`                    // Number of responses with the video segment or video manifest MIME type (i.e., application/x-mpegurl, application/vnd.apple.mpegurl, application/f4m, application/dash+xml, application/vnd.ms-sstr+xml, ideo/mp2t, audio/aac, video/f4f, video/x-flv, video/mp4, audio/mp4).
	WAFBlocked                *uint64     `mapstructure:"waf_blocked"`              // Number of requests that triggered a WAF rule and were blocked.
	WAFLogged                 *uint64     `mapstructure:"waf_logged"`               // Number of requests that triggered a WAF rule and were logged.
	WAFPassed                 *uint64     `mapstructure:"waf_passed"`               // Number of requests that triggered a WAF rule and were passed.
}

// GetStatsInput is an input to the GetStats function.
// Stats can be filtered by a Service ID, an individual stats field,
// time range (From and To), sampling rate (By) and/or Fastly region (Region)
// Allowed values for the fields are described at https://developer.fastly.com/reference/api/metrics-stats/
type GetStatsInput struct {
	// By is the duration of sample windows.
	By *string
	// Field is the name of the stats field.
	Field *string
	// From is the timestamp that defines the start of the window for which to fetch statistics, including the timestamp itself.
	From *string
	// Region limits query to a specific geographic region.
	Region *string
	// Service is the ID of the service.
	Service *string
	// To is the timestamp that defines the end of the window for which to fetch statistics.
	To *string
}

// GetStats retrieves the specified resource.
func (c *Client) GetStats(i *GetStatsInput) (*StatsResponse, error) {
	var resp any
	if err := c.GetStatsJSON(i, &resp); err != nil {
		return nil, err
	}

	var sr *StatsResponse
	if err := decodeMap(resp, &sr); err != nil {
		return nil, err
	}
	return sr, nil
}

// StatsFieldResponse is a response from the service stats/field API endpoint.
type StatsFieldResponse struct {
	Data    map[string][]*Stats `mapstructure:"data"`
	Message *string             `mapstructure:"msg"`
	Meta    map[string]string   `mapstructure:"meta"`
	Status  *string             `mapstructure:"status"`
}

// GetStatsField retrieves the specified resource.
func (c *Client) GetStatsField(i *GetStatsInput) (*StatsFieldResponse, error) {
	var resp any
	if err := c.GetStatsJSON(i, &resp); err != nil {
		return nil, err
	}

	var sr *StatsFieldResponse
	if err := decodeMap(resp, &sr); err != nil {
		return nil, err
	}
	return sr, nil
}

// GetStatsJSON fetches stats and decodes the response directly to the JSON struct dst.
func (c *Client) GetStatsJSON(i *GetStatsInput, dst any) error {
	p := "/stats"
	if i.Service != nil {
		p = fmt.Sprintf("%s/service/%s", p, *i.Service)
	}
	if i.Field != nil {
		p = fmt.Sprintf("%s/field/%s", p, *i.Field)
	}

	ro := &RequestOptions{
		Params: map[string]string{},
	}
	if i.By != nil {
		ro.Params["by"] = *i.By
	}
	if i.From != nil {
		ro.Params["from"] = *i.From
	}
	if i.Region != nil {
		ro.Params["region"] = *i.Region
	}
	if i.To != nil {
		ro.Params["to"] = *i.To
	}

	resp, err := c.Get(p, ro)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dst)
}

// Usage represents usage data of a single service or region.
type Usage struct {
	Bandwidth       *uint64 `mapstructure:"bandwidth"`
	Requests        *uint64 `mapstructure:"requests"`
	ComputeRequests *uint64 `mapstructure:"compute_requests"`
}

// UsageResponse is a response from the account usage API endpoint.
type UsageResponse struct {
	Data    *RegionsUsage     `mapstructure:"data"`
	Message *string           `mapstructure:"msg"`
	Meta    map[string]string `mapstructure:"meta"`
	Status  *string           `mapstructure:"status"`
}

// RegionsUsage is a list of aggregated usage data by Fastly's region.
type RegionsUsage map[string]*Usage

// GetUsageInput is used as an input to the GetUsage function
// Value for the input are described at https://developer.fastly.com/reference/api/metrics-stats/
type GetUsageInput struct {
	// By is the duration of sample windows.
	By *string
	// From is the timestamp that defines the start of the window for which to fetch statistics, including the timestamp itself.
	From *string
	// Region limits query to a specific geographic region.
	Region *string
	// To is the timestamp that defines the end of the window for which to fetch statistics.
	To *string
}

// GetUsage returns usage information aggregated across all Fastly services and grouped by region.
func (c *Client) GetUsage(i *GetUsageInput) (*UsageResponse, error) {
	ro := &RequestOptions{
		Params: map[string]string{},
	}
	if i.By != nil {
		ro.Params["by"] = *i.By
	}
	if i.From != nil {
		ro.Params["from"] = *i.From
	}
	if i.Region != nil {
		ro.Params["region"] = *i.Region
	}
	if i.To != nil {
		ro.Params["to"] = *i.To
	}

	resp, err := c.Get("/stats/usage", ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr *UsageResponse
	if err := decodeBodyMap(resp.Body, &sr); err != nil {
		return nil, err
	}

	return sr, nil
}

// UsageByServiceResponse is a response from the account usage API endpoint.
type UsageByServiceResponse struct {
	Data    *ServicesByRegionsUsage `mapstructure:"data"`
	Message *string                 `mapstructure:"msg"`
	Meta    map[string]string       `mapstructure:"meta"`
	Status  *string                 `mapstructure:"status"`
}

// ServicesUsage is a list of usage data by a service.
type ServicesUsage map[string]*Usage

// ServicesByRegionsUsage is a list of ServicesUsage by Fastly's region.
type ServicesByRegionsUsage map[string]*ServicesUsage

// GetUsageByService returns usage information aggregated by service and
// grouped by service and region.
func (c *Client) GetUsageByService(i *GetUsageInput) (*UsageByServiceResponse, error) {
	ro := &RequestOptions{
		Params: map[string]string{},
	}
	if i.By != nil {
		ro.Params["by"] = *i.By
	}
	if i.From != nil {
		ro.Params["from"] = *i.From
	}
	if i.Region != nil {
		ro.Params["region"] = *i.Region
	}
	if i.To != nil {
		ro.Params["to"] = *i.To
	}

	resp, err := c.Get("/stats/usage_by_service", ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr *UsageByServiceResponse
	if err := decodeBodyMap(resp.Body, &sr); err != nil {
		return nil, err
	}

	return sr, nil
}

// RegionsResponse is a response from Fastly regions API endpoint.
type RegionsResponse struct {
	Data    []string          `mapstructure:"data"`
	Message *string           `mapstructure:"msg"`
	Meta    map[string]string `mapstructure:"meta"`
	Status  *string           `mapstructure:"status"`
}

// GetRegions returns a list of Fastly regions.
func (c *Client) GetRegions() (*RegionsResponse, error) {
	resp, err := c.Get("stats/regions", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rr *RegionsResponse
	if err := decodeBodyMap(resp.Body, &rr); err != nil {
		return nil, err
	}

	return rr, nil
}
