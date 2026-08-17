package fastly

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

// LogInsightsVisualization represents a visualization supported by the Log Insights API.
type LogInsightsVisualization string

const (
	// LogInsightsVisualizationTopURLByBandwidth returns URLs using the most bandwidth.
	LogInsightsVisualizationTopURLByBandwidth LogInsightsVisualization = "top-url-by-bandwidth"
	// LogInsightsVisualizationBottomURLByCacheHitRatio returns URLs with the lowest cache hit ratio.
	LogInsightsVisualizationBottomURLByCacheHitRatio LogInsightsVisualization = "bottom-url-by-cache-hit-ratio"
	// LogInsightsVisualizationTopURLByCacheHitRatio returns URLs with the highest cache hit ratio.
	LogInsightsVisualizationTopURLByCacheHitRatio LogInsightsVisualization = "top-url-by-cache-hit-ratio"
	// LogInsightsVisualizationCountryStatistics returns country statistics.
	LogInsightsVisualizationCountryStatistics LogInsightsVisualization = "country-statistics"
	// LogInsightsVisualizationTopURLByDurationSum returns URLs with the highest cumulative response duration.
	LogInsightsVisualizationTopURLByDurationSum LogInsightsVisualization = "top-url-by-duration-sum"
	// LogInsightsVisualizationTop4XXURLs returns URLs with the most 4xx responses.
	LogInsightsVisualizationTop4XXURLs LogInsightsVisualization = "top-4xx-urls"
	// LogInsightsVisualizationTop5XXURLs returns URLs with the most 5xx responses.
	LogInsightsVisualizationTop5XXURLs LogInsightsVisualization = "top-5xx-urls"
	// LogInsightsVisualizationTopURLByMisses returns URLs with the most cache misses.
	LogInsightsVisualizationTopURLByMisses LogInsightsVisualization = "top-url-by-misses"
	// LogInsightsVisualizationTopURLByRequests returns URLs with the most requests.
	LogInsightsVisualizationTopURLByRequests LogInsightsVisualization = "top-url-by-requests"
	// LogInsightsVisualizationTopBrowserByRequests returns browsers generating the most requests.
	LogInsightsVisualizationTopBrowserByRequests LogInsightsVisualization = "top-browser-by-requests"
	// LogInsightsVisualizationTopContentTypeByRequests returns content types with the most requests.
	LogInsightsVisualizationTopContentTypeByRequests LogInsightsVisualization = "top-content-type-by-requests"
	// LogInsightsVisualizationTopDeviceByRequests returns device types generating the most requests.
	LogInsightsVisualizationTopDeviceByRequests LogInsightsVisualization = "top-device-by-requests"
	// LogInsightsVisualizationTopOSByRequests returns operating systems generating the most requests.
	LogInsightsVisualizationTopOSByRequests LogInsightsVisualization = "top-os-by-requests"
	// LogInsightsVisualizationResponseStatusCodes returns response status code statistics.
	LogInsightsVisualizationResponseStatusCodes LogInsightsVisualization = "response-status-codes"
	// LogInsightsVisualizationTop503Responses returns the most common 503 responses.
	LogInsightsVisualizationTop503Responses LogInsightsVisualization = "top-503-responses"
)

// LogInsightsVisualizations is a list of supported Log Insights visualizations.
var LogInsightsVisualizations = []LogInsightsVisualization{
	LogInsightsVisualizationTopURLByBandwidth,
	LogInsightsVisualizationBottomURLByCacheHitRatio,
	LogInsightsVisualizationTopURLByCacheHitRatio,
	LogInsightsVisualizationCountryStatistics,
	LogInsightsVisualizationTopURLByDurationSum,
	LogInsightsVisualizationTop4XXURLs,
	LogInsightsVisualizationTop5XXURLs,
	LogInsightsVisualizationTopURLByMisses,
	LogInsightsVisualizationTopURLByRequests,
	LogInsightsVisualizationTopBrowserByRequests,
	LogInsightsVisualizationTopContentTypeByRequests,
	LogInsightsVisualizationTopDeviceByRequests,
	LogInsightsVisualizationTopOSByRequests,
	LogInsightsVisualizationResponseStatusCodes,
	LogInsightsVisualizationTop503Responses,
}

// LogInsightsDimensions represents the dimensions returned by the Log Insights API.
type LogInsightsDimensions struct {
	Browser        *string `json:"browser"`
	BrowserVersion *string `json:"browser_version"`
	ContentType    *string `json:"content_type"`
	Country        *string `json:"country"`
	Device         *string `json:"device"`
	OS             *string `json:"os"`
	Region         *string `json:"region"`
	Response       *string `json:"response"`
	StatusCode     *string `json:"status-code"`
	URL            *string `json:"url"`
}

// LogInsightsValue represents the metrics returned for a Log Insights dimension.
type LogInsightsValue struct {
	AverageBandwidthBytes  *float64 `json:"average_bandwidth_bytes"`
	AverageResponseTime    *float64 `json:"average_response_time"`
	BandwidthPercentage    *float64 `json:"bandwidth_percentage"`
	CacheHitRatio          *float64 `json:"cache_hit_ratio"`
	CountryCHR             *float64 `json:"country_chr"`
	CountryErrorRate       *float64 `json:"country_error_rate"`
	CountryRequestRate     *float64 `json:"country_request_rate"`
	MissRate               *float64 `json:"miss_rate"`
	P95ResponseTime        *float64 `json:"p95_response_time"`
	Rate                   *float64 `json:"rate"`
	Rate503PerURL          *float64 `json:"503_rate_per_url"`
	RatePerStatus          *float64 `json:"rate_per_status"`
	RatePerURL             *float64 `json:"rate_per_url"`
	RegionCHR              *float64 `json:"region_chr"`
	RegionErrorRate        *float64 `json:"region_error_rate"`
	RequestPercentage      *float64 `json:"request_percentage"`
	ResponseTimePercentage *float64 `json:"response_time_percentage"`
}

// LogInsightsData represents a dimension and its associated values.
type LogInsightsData struct {
	Dimensions *LogInsightsDimensions `json:"dimensions"`
	Values     []*LogInsightsValue    `json:"values"`
}

// LogInsightsFilters represents the filters echoed in a Log Insights response.
type LogInsightsFilters struct {
	Domain           *string  `json:"domain"`
	DomainExactMatch *bool    `json:"domain_exact_match"`
	End              *string  `json:"end"`
	Limit            *int     `json:"limit"`
	POPs             []string `json:"pops"`
	ServiceID        *string  `json:"service_id"`
	Start            *string  `json:"start"`
}

// LogInsightsMeta represents metadata returned by the Log Insights API.
type LogInsightsMeta struct {
	Filters *LogInsightsFilters `json:"filters"`
}

// LogInsightsResponse represents the response from the Log Insights API.
type LogInsightsResponse struct {
	Data []*LogInsightsData `json:"data"`
	Meta *LogInsightsMeta   `json:"meta"`
}

// GetLogInsightsInput is used as input to the GetLogInsights function.
type GetLogInsightsInput struct {
	// Domain limits data to the specified request domain.
	Domain *string
	// DomainExactMatch determines whether Domain is treated as an exact match. The API defaults to true.
	DomainExactMatch *bool
	// End specifies the exclusive end time in RFC3339 format (required).
	End string
	// Limit is the maximum number of rows to return. The API defaults to 10 and limits the value to 100.
	Limit *int
	// POPs limits data to the supplied Fastly POP codes.
	POPs []string
	// ServiceID is the ID of the service for which insights should be returned (required).
	ServiceID string
	// Start specifies the inclusive start time in RFC3339 format (required).
	Start string
	// Visualization specifies the Log Insights visualization to return (required).
	Visualization LogInsightsVisualization
}

// GetLogInsights retrieves statistics from sampled log records.
func (c *Client) GetLogInsights(ctx context.Context, i *GetLogInsightsInput) (*LogInsightsResponse, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.Start == "" {
		return nil, ErrMissingStart
	}
	if i.End == "" {
		return nil, ErrMissingEnd
	}
	if i.Visualization == "" {
		return nil, ErrMissingVisualization
	}

	requestOptions := CreateRequestOptions()
	requestOptions.Params["service_id"] = i.ServiceID
	requestOptions.Params["start"] = i.Start
	requestOptions.Params["end"] = i.End
	requestOptions.Params["visualization"] = string(i.Visualization)

	if i.Domain != nil {
		requestOptions.Params["domain"] = *i.Domain
	}
	if i.DomainExactMatch != nil {
		requestOptions.Params["domain_exact_match"] = strconv.FormatBool(*i.DomainExactMatch)
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if len(i.POPs) > 0 {
		requestOptions.Params["pops"] = strings.Join(i.POPs, ",")
	}

	resp, err := c.Get(ctx, ToSafeURL("observability", "log-insights"), requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *LogInsightsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
