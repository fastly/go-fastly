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
	Errors []LoggingEndpointError
}

type LoggingEndpointError struct {
	SequenceNumber uint64 `json:"s"`
	Timestamp      uint64 `json:"t"`
	Stream         string `json:"o"`
	RequestID      string `json:"r"`
	Message        string `json:"m"`
	Endpoint       string `json:"e"`
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
