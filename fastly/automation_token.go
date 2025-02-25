package fastly

import (
	"fmt"
	"net/http"
	"time"
)

// AutomationTokenRole is used to match possible automation token roles.
type AutomationTokenRole string

const (
	// BillingRole allows view access to basic information about service configurations,
	// invoices, and account billing history.
	BillingRole AutomationTokenRole = "billing"
	// EngineerRole allows creating services and managing their configurations.
	EngineerRole AutomationTokenRole = "engineer"
	// UserRole allows view access to basic information about service configurations,
	// and controls.
	UserRole AutomationTokenRole = "user"
)

// AutomationTokenPaginator is used for pagination on AutomationToken endpoints.
// as they return JSONAPI data.
type AutomationTokenPaginator struct {
	Data []*AutomationToken           `mapstructure:"data"`
	Meta AutomationTokenPaginatorMeta `mapstructure:"meta"`
}

// AutomationTokenPaginatorMeta represents the metadata for an AutomationTokenPaginator.
type AutomationTokenPaginatorMeta struct {
	CurrentPage int `mapstructure:"current_page"`
	PerPage     int `mapstructure:"per_page"`
	RecordCount int `mapstructure:"record_count"`
	TotalPages  int `mapstructure:"total_pages"`
}

// AutomationToken represents an API token which allows non-human clients to
// authenticate requests to the Fastly API.
type AutomationToken struct {
	AccessToken *string              `mapstructure:"access_token"`
	CreatedAt   *time.Time           `mapstructure:"created_at"`
	ExpiresAt   *time.Time           `mapstructure:"expires_at"`
	IP          *string              `mapstructure:"ip"`
	LastUsedAt  *time.Time           `mapstructure:"last_used_at"`
	Name        *string              `mapstructure:"name"`
	Role        *AutomationTokenRole `mapstructure:"role"`
	Scope       *TokenScope          `mapstructure:"scope"`
	Services    []string             `mapstructure:"services"`
	TLSAccess   *bool                `mapstructure:"tls_access"`
	TokenID     *string              `mapstructure:"id"`
	UserID      *string              `mapstructure:"user_id"`
	CustomerID  *string              `mapstructure:"customer_id"`
}

// GetAutomationTokensInput is used as input to the GetAutomationTokens function.
type GetAutomationTokensInput struct {
	// Page is the current page.
	Page *int
	// PerPage is the number of records per page.
	PerPage *int
}

// GetAutomationTokens retrieves all resources.
func (c *Client) GetAutomationTokens(i *GetAutomationTokensInput) *ListPaginator[AutomationTokenPaginator] {
	input := ListOpts{}
	if i.Page != nil {
		input.Page = *i.Page
	}
	if i.PerPage != nil {
		input.PerPage = *i.PerPage
	}
	return NewPaginator[AutomationTokenPaginator](c, input, "/automation-tokens")
}

// ListAutomationTokens retrieves all resources.
func (c *Client) ListAutomationTokens() ([]*AutomationToken, error) {
	p := c.GetAutomationTokens(&GetAutomationTokensInput{})
	var results []*AutomationToken
	for p.HasNext() {
		data, err := p.GetNext()
		if err != nil {
			return nil, fmt.Errorf("failed to get next page (remaining: %d): %s", p.Remaining(), err)
		}

		for _, t := range data {
			results = append(results, t.Data...)
		}
	}
	return results, nil
}

// GetAutomationTokenInput is used as input to the GetAutomationToken function.
type GetAutomationTokenInput struct {
	// TokenID is an alphanumeric string identifying the token (required).
	TokenID string
}

// GetAutomationToken retrieves a specific resource by ID.
func (c *Client) GetAutomationToken(i *GetAutomationTokenInput) (*AutomationToken, error) {
	if i.TokenID == "" {
		return nil, ErrMissingTokenID
	}

	path := ToSafeURL("automation-tokens", i.TokenID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t *AutomationToken
	if err := DecodeBodyMap(resp.Body, &t); err != nil {
		return nil, err
	}
	return t, nil
}

// CreateAutomationTokenInput is used as input to the CreateAutomationToken function.
type CreateAutomationTokenInput struct {
	// ExpiresAt is a time-stamp (UTC) of when the token will expire.
	ExpiresAt *time.Time `json:"expires_at,omitempty" url:"expires_at,omitempty"`
	// Name is the name of the token.
	Name string `json:"name" url:"name,omitempty"`
	// Password is the password of the user the token is assigned to.
	Password *string `json:"-" url:"password,omitempty"`
	// Role is the role on the token (billing, engineer, user).
	Role AutomationTokenRole `json:"role" url:"role,omitempty"`
	// Scope is a space-delimited list of authorization scope (global, purge_select, purge_all, global).
	Scope *TokenScope `json:"scope,omitempty" url:"scope,omitempty"`
	// Services is a list of alphanumeric strings identifying services.
	// If no services are specified, the token will have access to all services on the account.
	Services []string `json:"services" url:"services,omitempty,brackets"`
	// Username is the email of the user the token is assigned to.
	Username *string `json:"-" url:"username,omitempty"`
	// TLSAccess indicates whether TLS access is enabled for the token.
	TLSAccess bool `json:"tls_access" url:"tls_access,omitempty"`
}

// CreateAutomationToken creates a new resource.
//
// Requires sudo capability for the token being used.
func (c *Client) CreateAutomationToken(i *CreateAutomationTokenInput) (*AutomationToken, error) {
	ignored, err := c.PostForm("/sudo", i, nil)
	if err != nil {
		return nil, err
	}
	defer ignored.Body.Close()

	resp, err := c.PostJSON("/automation-tokens", i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t *AutomationToken
	if err := DecodeBodyMap(resp.Body, &t); err != nil {
		return nil, err
	}
	return t, nil
}

// DeleteAutomationTokenInput is used as input to the DeleteAutomationToken function.
type DeleteAutomationTokenInput struct {
	// TokenID is an alphanumeric string identifying a token (required).
	TokenID string
}

// DeleteAutomationToken deletes the specified resource.
func (c *Client) DeleteAutomationToken(i *DeleteAutomationTokenInput) error {
	if i.TokenID == "" {
		return ErrMissingTokenID
	}

	path := ToSafeURL("tokens", i.TokenID)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return ErrNotOK
	}
	return nil
}
