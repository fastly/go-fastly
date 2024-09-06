package fastly

import (
	"net/http"
	"time"
)

// TokenScope is used to match possible authorization scopes.
type TokenScope string

const (
	// GlobalScope is the default scope covering all supported capabilities.
	GlobalScope TokenScope = "global"
	// PurgeSelectScope allows purging with surrogate key and URL, disallows purging with purge all.
	PurgeSelectScope TokenScope = "purge_select"
	// PurgeAllScope allows purging an entire service via purge_all.
	PurgeAllScope TokenScope = "purge_all"
	// GlobalReadScope allows read-only access to account information, configuration, and stats.
	GlobalReadScope TokenScope = "global:read"
)

// Token represents an API token which are used to authenticate requests to the
// Fastly API.
type Token struct {
	AccessToken *string     `mapstructure:"access_token"`
	CreatedAt   *time.Time  `mapstructure:"created_at"`
	ExpiresAt   *time.Time  `mapstructure:"expires_at"`
	IP          *string     `mapstructure:"ip"`
	LastUsedAt  *time.Time  `mapstructure:"last_used_at"`
	Name        *string     `mapstructure:"name"`
	Scope       *TokenScope `mapstructure:"scope"`
	Services    []string    `mapstructure:"services"`
	TokenID     *string     `mapstructure:"id"`
	UserID      *string     `mapstructure:"user_id"`
}

// ListTokensInput is used as input to the ListTokens function.
type ListTokensInput struct {
	// For backward compatibility.
}

// ListTokens retrieves all resources.
func (c *Client) ListTokens(_ *ListTokensInput) ([]*Token, error) {
	resp, err := c.Get("/tokens", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t []*Token
	if err := decodeBodyMap(resp.Body, &t); err != nil {
		return nil, err
	}
	return t, nil
}

// ListCustomerTokensInput is used as input to the ListCustomerTokens function.
type ListCustomerTokensInput struct {
	// CustomerID is an alphanumeric string identifying the customer (required).
	CustomerID string
}

// ListCustomerTokens retrieves all resources.
func (c *Client) ListCustomerTokens(i *ListCustomerTokensInput) ([]*Token, error) {
	if i.CustomerID == "" {
		return nil, ErrMissingCustomerID
	}

	path := ToSafeURL("customer", i.CustomerID, "tokens")

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t []*Token
	if err := decodeBodyMap(resp.Body, &t); err != nil {
		return nil, err
	}
	return t, nil
}

// GetTokenSelf retrieves the token information for the the access_token used
// used to authenticate the request.
//
// Returns a 401 if the token has expired and a 403 for invalid access token.
func (c *Client) GetTokenSelf() (*Token, error) {
	resp, err := c.Get("/tokens/self", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t *Token
	if err := decodeBodyMap(resp.Body, &t); err != nil {
		return nil, err
	}

	return t, nil
}

// CreateTokenInput is used as input to the Token function.
type CreateTokenInput struct {
	// ExpiresAt is a time-stamp (UTC) of when the token will expire
	ExpiresAt *time.Time `url:"expires_at,omitempty"`
	// Name is the name of the token.
	Name *string `url:"name,omitempty"`
	// Password is the token password.
	Password *string `url:"password,omitempty"`
	// Scope is a space-delimited list of authorization scope (global, purge_select, purge_all, global).
	Scope *TokenScope `url:"scope,omitempty"`
	// Services is a list of alphanumeric strings identifying services. If no services are specified, the token will have access to all services on the account.
	Services []string `url:"services,brackets,omitempty"`
	// Username is the email of the user the token is assigned to.
	Username *string `url:"username,omitempty"`
}

// CreateToken creates a new resource.
func (c *Client) CreateToken(i *CreateTokenInput) (*Token, error) {
	_, err := c.PostForm("/sudo", i, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.PostForm("/tokens", i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t *Token
	if err := decodeBodyMap(resp.Body, &t); err != nil {
		return nil, err
	}
	return t, nil
}

// DeleteTokenInput is used as input to the DeleteToken function.
type DeleteTokenInput struct {
	// TokenID is an alphanumeric string identifying a token (required).
	TokenID string
}

// DeleteToken deletes the specified resource.
func (c *Client) DeleteToken(i *DeleteTokenInput) error {
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

// DeleteTokenSelf deletes the specified resource.
func (c *Client) DeleteTokenSelf() error {
	resp, err := c.Delete("/tokens/self", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return ErrNotOK
	}
	return nil
}

// BatchDeleteTokensInput is used as input to BatchDeleteTokens.
type BatchDeleteTokensInput struct {
	// Tokens is a list of alphanumeric strings, each identifying a token.
	Tokens []*BatchToken
}

// BatchToken represents the JSONAPI data to be sent to the API.
// Reference: https://github.com/google/jsonapi#primary
type BatchToken struct {
	// ID is an alphanumeric string identifying a token.
	ID string `jsonapi:"primary,token,omitempty"`
}

// BatchDeleteTokens revokes multiple tokens.
func (c *Client) BatchDeleteTokens(i *BatchDeleteTokensInput) error {
	if len(i.Tokens) == 0 {
		return ErrMissingTokensValue
	}
	_, err := c.DeleteJSONAPIBulk("/tokens", i.Tokens, nil)
	return err
}
