package fastly

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// GetPrivateKeyInput is an input to the GetPrivateKey function.
// Allowed values for the fields are described at https://developer.fastly.com/reference/api/tls/platform/.
type GetPrivateKeyInput struct {
	// ID is an alphanumeric string identifying a private Key.
	ID string
}

// PrivateKey represents a private key is used to sign a Certificate.
type PrivateKey struct {
	CreatedAt     *time.Time `jsonapi:"attr,created_at,iso8601"`
	ID            string     `jsonapi:"primary,tls_private_key"`
	KeyLength     int        `jsonapi:"attr,key_length"`
	KeyType       string     `jsonapi:"attr,key_type"`
	Name          string     `jsonapi:"attr,name"`
	PublicKeySHA1 string     `jsonapi:"attr,public_key_sha1"`
	Replace       bool       `jsonapi:"attr,replace"`
}

// ListPrivateKeysInput is used as input to the ListPrivateKeys function.
type ListPrivateKeysInput struct {
	// FilterInUse is the returned keys to those without any matching TLS certificates.
	FilterInUse string
	// PageNumber is the page index for pagination.
	PageNumber int
	// PageSize is the number of keys per page.
	PageSize int
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListPrivateKeysInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"filter[in_use]": i.FilterInUse,
		"page[size]":     i.PageSize,
		"page[number]":   i.PageNumber,
	}

	for key, value := range pairings {
		switch v := value.(type) {
		case string:
			if v != "" {
				result[key] = v
			}
		case int:
			if v != 0 {
				result[key] = strconv.Itoa(v)
			}
		}
	}
	return result
}

// ListPrivateKeys retrieves all resources.
func (c *Client) ListPrivateKeys(i *ListPrivateKeysInput) ([]*PrivateKey, error) {
	path := "/tls/private_keys"
	filters := &RequestOptions{
		Params: i.formatFilters(),
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the filters don't work
		},
	}

	resp, err := c.Get(path, filters)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(PrivateKey)))
	if err != nil {
		return nil, err
	}

	ppk := make([]*PrivateKey, len(data))
	for i := range data {
		typed, ok := data[i].(*PrivateKey)
		if !ok {
			return nil, fmt.Errorf("got back a non-PrivateKey response")
		}
		ppk[i] = typed
	}

	return ppk, nil
}

// GetPrivateKey retrieves the specified resource.
func (c *Client) GetPrivateKey(i *GetPrivateKeyInput) (*PrivateKey, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("tls", "private_keys", i.ID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ppk PrivateKey
	if err := jsonapi.UnmarshalPayload(resp.Body, &ppk); err != nil {
		return nil, err
	}

	return &ppk, nil
}

// CreatePrivateKeyInput is used as input to the CreatePrivateKey function.
type CreatePrivateKeyInput struct {
	// Key is the contents of the private key. Must be a PEM-formatted key.
	Key string `jsonapi:"attr,key,omitempty"`
	// Name is a customizable name for your private key.
	Name string `jsonapi:"attr,name,omitempty"`
}

// CreatePrivateKey creates a new resource.
func (c *Client) CreatePrivateKey(i *CreatePrivateKeyInput) (*PrivateKey, error) {
	path := "/tls/private_keys"

	if i.Key == "" {
		return nil, ErrMissingKey
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	resp, err := c.PostJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ppk PrivateKey
	if err := jsonapi.UnmarshalPayload(resp.Body, &ppk); err != nil {
		return nil, err
	}

	return &ppk, nil
}

// DeletePrivateKeyInput used for deleting a private key.
type DeletePrivateKeyInput struct {
	// ID is an alphanumeric string identifying a private Key.
	ID string
}

// DeletePrivateKey deletes the specified resource.
func (c *Client) DeletePrivateKey(i *DeletePrivateKeyInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := ToSafeURL("tls", "private_keys", i.ID)

	_, err := c.Delete(path, nil)
	return err
}
