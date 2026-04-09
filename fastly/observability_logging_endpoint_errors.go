package fastly

import (
	"bufio"
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

type LoggingEndpointErrorsInput struct {
	// From is a unix-formatted timestamp to start the log stream from. Required for the initial request only.
	// If not used in conjunction with the `to` parameter, the range of the query will be a 10-second bucket.
	From *uint64
	// To is a unix-formatted timestamp to end the log stream. The maximum range between `from` and `to` is 1 hour;
	// requests exceeding this limit will return an error.
	To *uint64
	// Filter is a comma-separated list of logging endpoint names to filter the error stream.
	Filter []string
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string
}

type LoggingEndpointErrorsResponse struct {
	Errors   []LoggingEndpointError
	NextLink string
	PrevLink string
}

type LoggingEndpointError struct {
	SequenceNumber uint64 `json:"sequence_number"`
	Timestamp      uint64 `json:"error_time_us"`
	Stream         string `json:"stream"`
	Message        string `json:"message"`
	Endpoint       string `json:"endpoint"`
	Details        string `json:"details"`
}

func (c *Client) GetLoggingEndpointErrors(ctx context.Context, i *LoggingEndpointErrorsInput) (*LoggingEndpointErrorsResponse, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("observability", "service", i.ServiceID, "logging", "errors")
	requestOptions := CreateRequestOptions()

	if i.From != nil {
		requestOptions.Params["from"] = strconv.FormatUint(*i.From, 10)
	}
	if i.To != nil {
		requestOptions.Params["to"] = strconv.FormatUint(*i.To, 10)
	}
	if len(i.Filter) > 0 {
		requestOptions.Params["filter[endpoint]"] = strings.Join(i.Filter, ",")
	}

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result LoggingEndpointErrorsResponse

	// Parse Link header for pagination
	if linkHeader := resp.Header.Get("Link"); linkHeader != "" {
		result.NextLink, result.PrevLink = parseLinkHeader(linkHeader)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var errorLog LoggingEndpointError
		if err := json.Unmarshal(scanner.Bytes(), &errorLog); err != nil {
			return nil, err
		}
		result.Errors = append(result.Errors, errorLog)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

// parseLinkHeader parses the Link header and extracts the 'from' parameter from next and prev URLs.
// Format: </path?from=123>; rel="next", </path?from=456>; rel="prev"
func parseLinkHeader(header string) (next, prev string) {
	links := strings.Split(header, ",")
	for _, link := range links {
		link = strings.TrimSpace(link)
		parts := strings.Split(link, ";")
		if len(parts) != 2 {
			continue
		}

		url := strings.Trim(strings.TrimSpace(parts[0]), "<>")
		rel := strings.TrimSpace(parts[1])

		// Extract 'from' parameter from URL
		fromValue := extractFromParam(url)

		if strings.Contains(rel, `rel="next"`) {
			next = fromValue
		} else if strings.Contains(rel, `rel="prev"`) {
			prev = fromValue
		}
	}
	return next, prev
}

// extractFromParam extracts the 'from' parameter value from a URL or URL-encoded string.
func extractFromParam(url string) string {
	// Handle both encoded (%3Ffrom=) and unencoded (?from=) formats
	if idx := strings.Index(url, "from="); idx != -1 {
		fromStr := url[idx+5:]
		// Take until the next & or end of string
		if endIdx := strings.Index(fromStr, "&"); endIdx != -1 {
			fromStr = fromStr[:endIdx]
		}
		return fromStr
	}
	return ""
}
