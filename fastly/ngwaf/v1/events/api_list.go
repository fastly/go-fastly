package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v10/fastly"
)

// ListInput specifies the information needed for the List() function to perform
// the operation.
type ListInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
	// From represents the start of a date-time range, expressed in RFC 3339 format (required).
	From *string `url:"from"`
	// IP filters the list of events based on IP.
	IP *string `url:"ip"`
	// Limit is the limit on how many results are returned. [Default 100]
	Limit *int `url:"limit"`
	// Page is the page number of the collection to request. [Default 0]
	Page *int `url:"page"`
	// Signal filters the list of events based on signal.
	Signal *string `url:"signal"`
	// Status filters the list of events based on status. Must be one of `active` or `expired`.
	Status *string `url:"status"`
	// To represents the end of a date-time range and must be older than from, expressed in RFC 3339 format.
	To *string `url:"to"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `url:"-"`
}

// Get retrieves the specified workspace.
func List(c *fastly.Client, i *ListInput) (*Events, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.From == nil {
		return nil, fastly.ErrMissingFrom
	}
	requestOptions := fastly.CreateRequestOptions(i.Context)
	requestOptions.Params["from"] = *i.From
	if i.IP != nil {
		requestOptions.Params["ip"] = *i.IP
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Page != nil {
		requestOptions.Params["page"] = strconv.Itoa(*i.Page)
	}
	if i.Signal != nil {
		requestOptions.Params["signal"] = *i.Signal
	}
	if i.To != nil {
		requestOptions.Params["to"] = *i.To
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "events")

	resp, err := c.Get(path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var events *Events
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return events, nil
}
