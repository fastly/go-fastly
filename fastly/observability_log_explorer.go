package fastly

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// LogExplorerFilterField represents a field that can be filtered by the Log Explorer API.
type LogExplorerFilterField string

const (
	// LogExplorerFilterFieldDomain filters by request domain.
	LogExplorerFilterFieldDomain LogExplorerFilterField = "domain"
	// LogExplorerFilterFieldRequestPath filters by request path.
	LogExplorerFilterFieldRequestPath LogExplorerFilterField = "request_path"
	// LogExplorerFilterFieldFastlyPOP filters by Fastly POP.
	LogExplorerFilterFieldFastlyPOP LogExplorerFilterField = "fastly_pop"
	// LogExplorerFilterFieldResponseTime filters by response time.
	LogExplorerFilterFieldResponseTime LogExplorerFilterField = "response_time"
	// LogExplorerFilterFieldResponseStatus filters by HTTP response status.
	LogExplorerFilterFieldResponseStatus LogExplorerFilterField = "response_status"
	// LogExplorerFilterFieldFastlyIsShield filters by whether the request was handled by a shield POP.
	LogExplorerFilterFieldFastlyIsShield LogExplorerFilterField = "fastly_is_shield"
	// LogExplorerFilterFieldFastlyIsEdge filters by whether the request was handled by an edge POP.
	LogExplorerFilterFieldFastlyIsEdge LogExplorerFilterField = "fastly_is_edge"
	// LogExplorerFilterFieldClientOSName filters by client operating system name.
	LogExplorerFilterFieldClientOSName LogExplorerFilterField = "client_os_name"
	// LogExplorerFilterFieldClientDeviceType filters by client device type.
	LogExplorerFilterFieldClientDeviceType LogExplorerFilterField = "client_device_type"
	// LogExplorerFilterFieldClientBrowserName filters by client browser name.
	LogExplorerFilterFieldClientBrowserName LogExplorerFilterField = "client_browser_name"
	// LogExplorerFilterFieldFastlyIsCacheHit filters by whether the request was fulfilled from cache.
	LogExplorerFilterFieldFastlyIsCacheHit LogExplorerFilterField = "fastly_is_cache_hit"
)

// LogExplorerFilterFields is a list of supported Log Explorer filter fields.
var LogExplorerFilterFields = []LogExplorerFilterField{
	LogExplorerFilterFieldDomain,
	LogExplorerFilterFieldRequestPath,
	LogExplorerFilterFieldFastlyPOP,
	LogExplorerFilterFieldResponseTime,
	LogExplorerFilterFieldResponseStatus,
	LogExplorerFilterFieldFastlyIsShield,
	LogExplorerFilterFieldFastlyIsEdge,
	LogExplorerFilterFieldClientOSName,
	LogExplorerFilterFieldClientDeviceType,
	LogExplorerFilterFieldClientBrowserName,
	LogExplorerFilterFieldFastlyIsCacheHit,
}

// LogExplorerFilterOperator represents a comparison operator supported by the Log Explorer API.
type LogExplorerFilterOperator string

const (
	// LogExplorerFilterOperatorEq performs an equality comparison.
	LogExplorerFilterOperatorEq LogExplorerFilterOperator = "eq"
	// LogExplorerFilterOperatorEndsWith performs a string suffix comparison.
	LogExplorerFilterOperatorEndsWith LogExplorerFilterOperator = "ends-with"
	// LogExplorerFilterOperatorIn matches values in a list.
	LogExplorerFilterOperatorIn LogExplorerFilterOperator = "in"
	// LogExplorerFilterOperatorNotIn excludes values in a list.
	LogExplorerFilterOperatorNotIn LogExplorerFilterOperator = "not_in"
	// LogExplorerFilterOperatorGT performs a greater-than comparison.
	LogExplorerFilterOperatorGT LogExplorerFilterOperator = "gt"
	// LogExplorerFilterOperatorGTE performs a greater-than-or-equal comparison.
	LogExplorerFilterOperatorGTE LogExplorerFilterOperator = "gte"
	// LogExplorerFilterOperatorLT performs a less-than comparison.
	LogExplorerFilterOperatorLT LogExplorerFilterOperator = "lt"
	// LogExplorerFilterOperatorLTE performs a less-than-or-equal comparison.
	LogExplorerFilterOperatorLTE LogExplorerFilterOperator = "lte"
)

// LogExplorerFilterOperators is a list of supported Log Explorer filter operators.
var LogExplorerFilterOperators = []LogExplorerFilterOperator{
	LogExplorerFilterOperatorEq,
	LogExplorerFilterOperatorEndsWith,
	LogExplorerFilterOperatorIn,
	LogExplorerFilterOperatorNotIn,
	LogExplorerFilterOperatorGT,
	LogExplorerFilterOperatorGTE,
	LogExplorerFilterOperatorLT,
	LogExplorerFilterOperatorLTE,
}

// LogExplorerFilter represents a filter supplied to the Log Explorer API.
type LogExplorerFilter struct {
	// Field is the log field to which the filter is applied.
	Field LogExplorerFilterField
	// Operator is the comparison operator used for the filter.
	Operator LogExplorerFilterOperator
	// Value is the value to compare against.
	Value string
}

// LogRecord represents a sampled request log returned by the Log Explorer API.
type LogRecord struct {
	ClientASNumber        *uint64  `json:"client_as_number"`
	ClientBrowserName     *string  `json:"client_browser_name"`
	ClientBrowserVersion  *string  `json:"client_browser_version"`
	ClientCountryCode     *string  `json:"client_country_code"`
	ClientDeviceType      *string  `json:"client_device_type"`
	ClientOSName          *string  `json:"client_os_name"`
	ClientOSVersion       *string  `json:"client_os_version"`
	ClientRegion          *string  `json:"client_region"`
	CustomerID            *string  `json:"customer_id"`
	FastlyPOP             *string  `json:"fastly_pop"`
	IsCacheHit            *bool    `json:"is_cache_hit"`
	IsEdge                *bool    `json:"is_edge"`
	IsShield              *bool    `json:"is_shield"`
	OriginHost            *string  `json:"origin_host"`
	RequestHost           *string  `json:"request_host"`
	RequestMethod         *string  `json:"request_method"`
	RequestPath           *string  `json:"request_path"`
	RequestProtocol       *string  `json:"request_protocol"`
	ResponseBytesBody     *uint64  `json:"response_bytes_body"`
	ResponseBytesHeader   *uint64  `json:"response_bytes_header"`
	ResponseContentLength *uint64  `json:"response_content_length"`
	ResponseContentType   *string  `json:"response_content_type"`
	ResponseReason        *string  `json:"response_reason"`
	ResponseState         *string  `json:"response_state"`
	ResponseStatus        *int     `json:"response_status"`
	ResponseTime          *float64 `json:"response_time"`
	ResponseXCache        *string  `json:"response_x_cache"`
	ServiceID             *string  `json:"service_id"`
	Timestamp             *string  `json:"timestamp"`
}

// LogExplorerAppliedFilter represents a filter echoed by the Log Explorer API.
type LogExplorerAppliedFilter struct {
	Field    *LogExplorerFilterField    `json:"field"`
	Operator *LogExplorerFilterOperator `json:"operator"`
	Value    any                        `json:"value"`
}

// LogExplorerFilters represents the filters echoed in a Log Explorer response.
type LogExplorerFilters struct {
	End          *string                     `json:"end"`
	FieldFilters []*LogExplorerAppliedFilter `json:"field_filters"`
	Limit        *int                        `json:"limit"`
	ServiceID    *string                     `json:"service_id"`
	Start        *string                     `json:"start"`
}

// LogExplorerMeta represents metadata returned by the Log Explorer API.
type LogExplorerMeta struct {
	Filters    *LogExplorerFilters `json:"filters"`
	NextCursor *string             `json:"next_cursor"`
}

// LogRecordsResponse represents the response from the Log Explorer API.
type LogRecordsResponse struct {
	Data []*LogRecord     `json:"data"`
	Meta *LogExplorerMeta `json:"meta"`
}

// GetLogRecordsInput is used as input to the GetLogRecords function.
type GetLogRecordsInput struct {
	// End specifies the exclusive end time in RFC3339 format (required).
	End string
	// Filters limits returned records to logs matching all supplied filters.
	Filters []LogExplorerFilter
	// Limit is the maximum number of rows to return. The API defaults to 10 and limits the value to 100.
	Limit *int
	// NextCursor is the cursor returned by a previous request for the next page.
	NextCursor *string
	// ServiceID is the ID of the service for which log records should be returned (required).
	ServiceID string
	// Start specifies the inclusive start time in RFC3339 format (required).
	Start string
}

// GetLogRecords retrieves sampled log records from the Log Explorer API.
func (c *Client) GetLogRecords(ctx context.Context, i *GetLogRecordsInput) (*LogRecordsResponse, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.Start == "" {
		return nil, ErrMissingStart
	}
	if i.End == "" {
		return nil, ErrMissingEnd
	}

	requestOptions := CreateRequestOptions()
	requestOptions.Params["service_id"] = i.ServiceID
	requestOptions.Params["start"] = i.Start
	requestOptions.Params["end"] = i.End

	for _, filter := range i.Filters {
		key := fmt.Sprintf("filter[%s][%s]", filter.Field, filter.Operator)
		requestOptions.Params[key] = filter.Value
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.NextCursor != nil {
		requestOptions.Params["next_cursor"] = *i.NextCursor
	}

	resp, err := c.Get(ctx, ToSafeURL("observability", "log-explorer"), requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *LogRecordsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
