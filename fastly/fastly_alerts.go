package fastly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// AlertDefinition holds the configuration for one alert.
type AlertDefinition struct {
	ID                 string              `json:"id"`
	ServiceID          string              `json:"service_id"`
	ServiceCustomerID  string              `json:"service_customer_id,omitempty"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Metric             string              `json:"metric"`
	Source             string              `json:"source"`
	Dimensions         map[string][]string `json:"dimensions"`
	EvaluationStrategy map[string]any      `json:"evaluation_strategy"`
	IntegrationIDs     []string            `json:"integration_ids"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	CreatedBy          string              `json:"created_by,omitempty"`
}

// AlertHistory describes the status of an alert definition over a time range.
type AlertHistory struct {
	ID           string          `json:"id"`
	DefinitionID string          `json:"definition_id"`
	Definition   AlertDefinition `json:"definition"`
	Start        time.Time       `json:"start"`
	End          time.Time       `json:"end"`
	Status       string          `json:"status"`
}

// AlertDefinitionsResponse is the response for an alert definitions query.
type AlertDefinitionsResponse struct {
	Data []AlertDefinition `json:"data"`
	Meta AlertsMeta        `json:"meta"`
}

// AlertHistoryResponse is the response for an alert history query.
type AlertHistoryResponse struct {
	Data []AlertHistory `json:"data"`
	Meta AlertsMeta     `json:"meta"`
}

// AlertsMeta holds metadata about an alerts query.
type AlertsMeta struct {
	NextCursor string `json:"next_cursor"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Total      int    `json:"total"`
}

// ListAlertDefinitionsInput is used as input to the ListAlertDefinitions function.
type ListAlertDefinitionsInput struct {
	// ServiceID filters definitions by service (optional).
	ServiceID *string
	// ServiceCustomerID is for Fastly Staff use only (optional).
	ServiceCustomerID *string
	// Name filters definitions by name substring (optional).
	Name *string
	// CreatedBy is for Fastly Staff use only (optional).
	CreatedBy *string
	// Cursor is the pagination cursor from a previous request's meta (optional).
	Cursor *string
	// Limit is the maximum number of items included in each response (optional).
	Limit *int
	// Sort is the field on which to sort definitions (optional).
	Sort *string
}

// ListAlertDefinitions retrieves filtered, paginated alert definitions.
func (c *Client) ListAlertDefinitions(i *ListAlertDefinitionsInput) (*AlertDefinitionsResponse, error) {
	p := "/alerts/definitions"

	ro := &RequestOptions{
		Params: map[string]string{},
	}
	if i.ServiceID != nil {
		ro.Params["service_id"] = *i.ServiceID
	}
	if i.ServiceCustomerID != nil {
		ro.Params["service_customer_id"] = *i.ServiceCustomerID
	}
	if i.Name != nil {
		ro.Params["name"] = *i.Name
	}
	if i.CreatedBy != nil {
		ro.Params["created_by"] = *i.CreatedBy
	}
	if i.Cursor != nil {
		ro.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		ro.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Sort != nil {
		ro.Params["sort"] = *i.Sort
	}

	resp, err := c.Get(p, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var adr *AlertDefinitionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&adr); err != nil {
		return nil, err
	}

	return adr, nil
}

// CreateAlertDefinitionInput is used as input to the CreateAlertDefinition function.
type CreateAlertDefinitionInput struct {
	// ServiceID is the ID of the service that the alert is monitoring (required).
	ServiceID string `json:"service_id"`
	// ServiceCustomerID is for Fastly staff use only (optional).
	ServiceCustomerID string `json:"service_customer_id,omitempty"`
	// Name is the summary text of the alert (required, limit 255).
	Name string `json:"name"`
	// Description is additional text included in an alert notification (optional, limit 4096).
	Description string `json:"description"`
	// Metric is the name of the metric being monitored for alert evaluation (required).
	Metric string `json:"metric"`
	// Source is the metric source (required). Options are: 'stats', 'origins', 'domains'.
	Source string `json:"source"`
	// Dimensions are a list of origins or domains that the alert is restricted to.
	Dimensions map[string][]string `json:"dimensions"`
	// EvaluationStrategy is the evaluation strategy for the alert (required).
	EvaluationStrategy map[string]any `json:"evaluation_strategy"`
	// IntegrationIDs are IDs of integrations that notifications will be sent to.
	IntegrationIDs []string `json:"integration_ids"`
}

// CreateAlertDefinition creates a new alert definition.
func (c *Client) CreateAlertDefinition(i *CreateAlertDefinitionInput) (*AlertDefinition, error) {
	resp, err := c.PostJSON("/alerts/definitions", i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ad *AlertDefinition
	if err := json.NewDecoder(resp.Body).Decode(&ad); err != nil {
		return nil, err
	}
	return ad, nil
}

// GetAlertDefinitionInput is used as input to the GetAlertDefinition function.
type GetAlertDefinitionInput struct {
	// ID of definition to fetch (required).
	ID string
}

// GetAlertDefinition retrieves a specified alert definition.
func (c *Client) GetAlertDefinition(i *GetAlertDefinitionInput) (*AlertDefinition, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/alerts/definitions/%s", i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ad *AlertDefinition
	if err := json.NewDecoder(resp.Body).Decode(&ad); err != nil {
		return nil, err
	}

	return ad, nil
}

// UpdateAlertDefinitionInput is used as input to the UpdateAlertDefinition function.
type UpdateAlertDefinitionInput struct {
	// ID of definition to update (required).
	ID string `json:"-"`
	// Name is the summary text of the alert (required, limit 255).
	Name string `json:"name"`
	// Description is additional text included in an alert notification (optional, limit 4096).
	Description string `json:"description"`
	// Metric is the name of the metric being monitored for alert evaluation (required).
	Metric string `json:"metric"`
	// Source is the metric source (required). Options are: 'stats', 'origins', 'domains'.
	Source string `json:"source"`
	// Dimensions are a list of origins or domains that the alert is restricted to.
	Dimensions map[string][]string `json:"dimensions"`
	// EvaluationStrategy is the evaluation strategy for the alert (required).
	EvaluationStrategy map[string]any `json:"evaluation_strategy"`
	// IntegrationIDs are IDs of integrations that notifications will be sent to.
	IntegrationIDs []string `json:"integration_ids"`
}

// UpdateAlertDefinition updates the specified alert definition.
func (c *Client) UpdateAlertDefinition(i *UpdateAlertDefinitionInput) (*AlertDefinition, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/alerts/definitions/%s", i.ID)
	resp, err := c.PutJSON(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ad *AlertDefinition
	if err := json.NewDecoder(resp.Body).Decode(&ad); err != nil {
		return nil, err
	}
	return ad, nil
}

// DeleteAlertDefinitionInput is used as input to the DeleteAlertDefinition function.
type DeleteAlertDefinitionInput struct {
	// ID of definition to delete (required).
	ID string
}

// DeleteAlertDefinition deletes the specified alert definition.
func (c *Client) DeleteAlertDefinition(i *DeleteAlertDefinitionInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/alerts/definitions/%s", i.ID)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}

// TestAlertDefinitionInput is used as input to the TestAlertDefinition function.
type TestAlertDefinitionInput struct {
	// Same as CreateAlertDefinitionInput
	CreateAlertDefinitionInput
}

// TestAlertDefinition validates alert definition and sends test notifications without creating.
func (c *Client) TestAlertDefinition(i *TestAlertDefinitionInput) error {
	resp, err := c.PostJSON("/alerts/definitions/test", i, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return NewHTTPError(resp)
	}

	return nil
}

// ListAlertHistoryInput is used as input to the ListAlertHistory function.
type ListAlertHistoryInput struct {
	// Status is the alert status.
	Status *string
	// ServiceID filters history by service (optional).
	ServiceID *string
	// ServiceCustomerID is for Fastly staff use only (optional).
	ServiceCustomerID *string
	// DefinitionID filters history by definition (optional).
	DefinitionID *string
	// Filter history to those with start or end on or after the provided timestamp (optional).
	After *string
	// Filter history to those with start or end on or before the provided timestamp (optional).
	Before *string
	// CreatedBy is for Fastly Staff use only (optional).
	CreatedBy *string
	// Cursor is the pagination cursor from a previous request's meta (optional).
	Cursor *string
	// Limit is the maximum number of items included in each response (optional).
	Limit *int
	// Sort is the field on which to sort definitions (optional).
	Sort *string
}

// ListAlertHistory retrieves filtered, paginated alert history records.
func (c *Client) ListAlertHistory(i *ListAlertHistoryInput) (*AlertHistoryResponse, error) {
	p := "/alerts/history"

	ro := &RequestOptions{
		Params: map[string]string{},
	}
	if i.Status != nil {
		ro.Params["status"] = *i.Status
	}
	if i.ServiceID != nil {
		ro.Params["service_id"] = *i.ServiceID
	}
	if i.ServiceCustomerID != nil {
		ro.Params["service_customer_id"] = *i.ServiceCustomerID
	}
	if i.DefinitionID != nil {
		ro.Params["definition_id"] = *i.DefinitionID
	}
	if i.After != nil {
		ro.Params["after"] = *i.After
	}
	if i.Before != nil {
		ro.Params["before"] = *i.Before
	}
	if i.CreatedBy != nil {
		ro.Params["created_by"] = *i.CreatedBy
	}
	if i.Cursor != nil {
		ro.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		ro.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Sort != nil {
		ro.Params["sort"] = *i.Sort
	}

	resp, err := c.Get(p, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ahr *AlertHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&ahr); err != nil {
		return nil, err
	}

	return ahr, nil
}
