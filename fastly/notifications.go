package fastly

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Integration holds the configuration for one integration.
type Integration struct {
	CreatedAt   *time.Time        `json:"created_at"`
	Description *string           `json:"description"`
	Config      map[string]string `json:"config"`
	ID          *string           `json:"id"`
	Name        *string           `json:"name"`
	Status      *string           `json:"status"`
	Type        *string           `json:"type"`
	UpdatedAt   *time.Time        `json:"updated_at"`
}

// SearchIntegrationsInput is used as input to the SearchIntegrations function.
type SearchIntegrationsInput struct {
	// Cursor is the pagination cursor from a previous request's meta.
	Cursor *string
	// Limit is the maximum number of items included in each response.
	Limit *int
	// Sort is the field on which to sort integrations.
	Sort *string
	// Type filters integrations by type.
	Type *string
}

// SearchIntegrationsResponse is the response for an integrations query.
type SearchIntegrationsResponse struct {
	Data []Integration     `json:"data"`
	Meta *IntegrationsMeta `json:"meta"`
}

// IntegrationsMeeta holds metadata about an integrations query.
type IntegrationsMeta struct {
	Limit      *int    `json:"limit"`
	NextCursor *string `json:"next_cursor"`
	Sort       *string `json:"sort"`
	Total      *int    `json:"total"`
	Type       *string `json:"type"`
}

// SearchIntegrations retrieves filtered, paginated integrations.
func (c *Client) SearchIntegrations(ctx context.Context, i *SearchIntegrationsInput) (*SearchIntegrationsResponse, error) {
	p := "/notifications/integrations"

	requestOptions := CreateRequestOptions()
	if i.Cursor != nil {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Sort != nil {
		requestOptions.Params["sort"] = *i.Sort
	}
	if i.Type != nil {
		requestOptions.Params["type"] = *i.Type
	}

	resp, err := c.Get(ctx, p, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sir *SearchIntegrationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&sir); err != nil {
		return nil, err
	}

	return sir, nil
}

// CreateIntegrationInput is used as input to the CreateIntegration function.
type CreateIntegrationInput struct {
	// Config is configuration specific to the integration type.
	Config map[string]string
	// Description is the user submitted description of the integration.
	Description *string
	// Name is the user submitted name of the integration.
	Name *string
	// Type is the type of integration.
	Type *string
}

// CreateIntegrationResponse is the response for creating a new integration.
type CreateIntegrationResponse struct {
	// ID of created integration.
	ID *string `json:"integration_id"`
}

// CreateIntegration creates a new integration.
func (c *Client) CreateIntegration(ctx context.Context, i *CreateIntegrationInput) (*CreateIntegrationResponse, error) {
	resp, err := c.PostJSON(ctx, "/notifications/integrations", i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cir *CreateIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&cir); err != nil {
		return nil, err
	}
	return cir, nil
}

// GetIntegrationInput is used as input to the GetIntegration function.
type GetIntegrationInput struct {
	// ID of integration to fetch (required).
	ID string
}

// GetIntegration retrieves a specified integration.
func (c *Client) GetIntegration(ctx context.Context, i *GetIntegrationInput) (*Integration, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("notifications", "integrations", i.ID)

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var integration *Integration
	if err := json.NewDecoder(resp.Body).Decode(&integration); err != nil {
		return nil, err
	}

	return integration, nil
}

// UpdateIntegrationInput is used as input to the UpdateIntegration function.
type UpdateIntegrationInput struct {
	// Config is configuration specific to the integration type.
	Config map[string]string
	// Description is the user submitted description of the integration.
	Description *string
	// ID of integration to update (required).
	ID string
	// Name is the user submitted name of the integration.
	Name *string
	// Type is the type of integration
	Type *string
}

// UpdateIntegration updates the specified integration.
func (c *Client) UpdateIntegration(ctx context.Context, i *UpdateIntegrationInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := ToSafeURL("notifications", "integrations", i.ID)

	resp, err := c.PatchJSON(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}

// DeleteIntegrationInput is used as input to the DeleteIntegration function.
type DeleteIntegrationInput struct {
	// ID of integration to delete (required).
	ID string
}

// DeleteIntegration deletes the specified integration.
func (c *Client) DeleteIntegration(ctx context.Context, i *DeleteIntegrationInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := ToSafeURL("notifications", "integrations", i.ID)

	resp, err := c.Delete(ctx, path, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}

// IntegrationType is an item in the response listing integration types.
type IntegrationType struct {
	Type         *string       `json:"type"`
	DisplayName  *string       `json:"display_name"`
	CustomFields []CustomField `json:"custom_fields"`
}

// CustomField describes a configuration required for a type of integration.
type CustomField struct {
	Name        *string `json:"name"`
	DisplayName *string `json:"display_name"`
	Format      *string `json:"format"`
}

// GetIntegrationTypes retrieves the supported integration types and what configuration they require.
func (c *Client) GetIntegrationTypes(ctx context.Context) (*[]IntegrationType, error) {
	path := "/notifications/integration-types"
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var its *[]IntegrationType
	if err := json.NewDecoder(resp.Body).Decode(&its); err != nil {
		return nil, err
	}
	return its, nil
}

// GetWebhookSigningKeyInput is used as input to the GetWebhookSigningKey function.
type GetWebhookSigningKeyInput struct {
	// IntegrationID is the ID of the webhook integration which signing key to get (required).
	IntegrationID string
}

// WebhookSigningKeyResponse is the response for getting or rotating a webhook payload signing key.
type WebhookSigningKeyResponse struct {
	SigningKey *string `json:"signingKey"`
}

// GetWebhookSigningKey retrieves the signing key for a webhook integration.
func (c *Client) GetWebhookSigningKey(ctx context.Context, i *GetWebhookSigningKeyInput) (*WebhookSigningKeyResponse, error) {
	if i.IntegrationID == "" {
		return nil, ErrMissingIntegrationID
	}

	path := ToSafeURL("notifications", "integrations", i.IntegrationID, "signingKey")

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wskr *WebhookSigningKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&wskr); err != nil {
		return nil, err
	}
	return wskr, nil
}

// RotateWebhookSigningKeyInput is used as input to the RotateWebhookSigningKey function.
type RotateWebhookSigningKeyInput struct {
	// IntegrationID is the ID of the webhook integration which signing key to rotate (required).
	IntegrationID string
}

// RotateWebhookSigningKey rotates the signing key for a webhook integration.
func (c *Client) RotateWebhookSigningKey(ctx context.Context, i *RotateWebhookSigningKeyInput) (*WebhookSigningKeyResponse, error) {
	if i.IntegrationID == "" {
		return nil, ErrMissingIntegrationID
	}

	path := ToSafeURL("notifications", "integrations", i.IntegrationID, "rotateSigningKey")

	resp, err := c.Post(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wskr *WebhookSigningKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&wskr); err != nil {
		return nil, err
	}
	return wskr, nil
}

// CreateMailinglistConfirmationInput is used as input to the CreateMailinglistConfirmation function.
type CreateMailinglistConfirmationInput struct {
	// Email is the mailinglist address.
	Email *string
}

// CreateMailinglistConfirmation sends a mailing list confirmation email.
func (c *Client) CreateMailinglistConfirmation(ctx context.Context, i *CreateMailinglistConfirmationInput) error {
	path := "/notifications/mailinglist-confirmations"
	resp, err := c.PostJSON(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}
