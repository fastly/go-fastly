package fastly

import (
	"context"
	"net/http"
)

// EdgeCheck represents an edge check response from the Fastly API.
type EdgeCheck struct {
	Hash         *string            `mapstructure:"hash"`
	Request      *EdgeCheckRequest  `mapstructure:"request"`
	Response     *EdgeCheckResponse `mapstructure:"response"`
	ResponseTime *float64           `mapstructure:"response_time"`
	Server       *string            `mapstructure:"server"`
}

// EdgeCheckRequest is the request part of an EdgeCheck response.
type EdgeCheckRequest struct {
	Headers *http.Header `mapstructure:"headers"`
	Method  *string      `mapstructure:"method"`
	URL     *string      `mapstructure:"url"`
}

// EdgeCheckResponse is the response part of an EdgeCheck response.
type EdgeCheckResponse struct {
	Headers *http.Header `mapstructure:"headers"`
	Status  *int         `mapstructure:"status"`
}

// EdgeCheckInput is used as input to the EdgeCheck function.
type EdgeCheckInput struct {
	// URL is the full URL (host and path) to check on all nodes.
	// If protocol is omitted, http will be assumed (required).
	URL string `url:"url,omitempty"`
}

// EdgeCheck queries the edge cache for all of Fastly's servers for the given
// URL.
func (c *Client) EdgeCheck(ctx context.Context, i *EdgeCheckInput) ([]*EdgeCheck, error) {
	var requestOptions RequestOptions
	if i != nil {
		requestOptions = CreateRequestOptions()
		requestOptions.Params["url"] = i.URL
	} else {
		requestOptions = CreateRequestOptions()
	}
	requestOptions.Parallel = true

	resp, err := c.Get(ctx, "/content/edge_check", requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e []*EdgeCheck
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}
