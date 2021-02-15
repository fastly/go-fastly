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
	ID string
}

// PrivateKey represents a private key is used to sign a Certificate.
type PrivateKey struct {
	ID            string     `jsonapi:"primary,tls_private_key"`
	Name          string     `jsonapi:"attr,name"`
	KeyLength     int        `jsonapi:"attr,key_length"`
	KeyType       string     `jsonapi:"attr,key_type"`
	PublicKeySHA1 string     `jsonapi:"attr,public_key_sha1"`
	CreatedAt     *time.Time `jsonapi:"attr,created_at,iso8601"`
	Replace       bool       `jsonapi:"attr,replace"`
}

// ListPrivateKeysInput is used as input to the ListPrivateKeys function.
type ListPrivateKeysInput struct {
	PageNumber  int    // The page index for pagination.
	PageSize    int    // The number of keys per page.
	FilterInUse string // Limit the returned keys to those without any matching TLS certificates.
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListPrivateKeysInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[in_use]": i.FilterInUse,
		"page[size]":     i.PageSize,
		"page[number]":   i.PageNumber,
	}

	for key, value := range pairings {
		switch t := reflect.TypeOf(value).String(); t {
		case "string":
			if value != "" {
				result[key] = value.(string)
			}
		case "int":
			if value != 0 {
				result[key] = strconv.Itoa(value.(int))
			}
		}
	}
	return result
}

// ListPrivateKeys list all TLS private keys.
func (c *Client) ListPrivateKeys(i *ListPrivateKeysInput) ([]*PrivateKey, error) {

	p := "/tls/private_keys"
	filters := &RequestOptions{
		Params: i.formatFilters(),
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the filters don't work
		},
	}

	r, err := c.Get(p, filters)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(PrivateKey)))
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

// GetPrivateKey show a TLS private key.
func (c *Client) GetPrivateKey(i *GetPrivateKeyInput) (*PrivateKey, error) {

	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := fmt.Sprintf("/tls/private_keys/%s", i.ID)

	r, err := c.Get(p, nil)
	if err != nil {
		return nil, err
	}

	var ppk PrivateKey
	if err := jsonapi.UnmarshalPayload(r.Body, &ppk); err != nil {
		return nil, err
	}

	return &ppk, nil
}

// CreatePrivateKeyInput is used as input to the CreatePrivateKey function.
type CreatePrivateKeyInput struct {
	Key  string `jsonapi:"attr,key,omitempty"`
	Name string `jsonapi:"attr,name,omitempty"`
}

// CreatePrivateKey create a TLS private key.
func (c *Client) CreatePrivateKey(i *CreatePrivateKeyInput) (*PrivateKey, error) {

	p := "/tls/private_keys"

	if i.Key == "" {
		return nil, ErrMissingKey
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	r, err := c.PostJSONAPI(p, i, nil)
	if err != nil {
		return nil, err
	}

	var ppk PrivateKey
	if err := jsonapi.UnmarshalPayload(r.Body, &ppk); err != nil {
		return nil, err
	}

	return &ppk, nil
}

// DeletePrivateKeyInput used for deleting a private key.
type DeletePrivateKeyInput struct {
	ID string
}

// DeletePrivateKey destroy a TLS private key. Only private keys not already matched to any certificates can be deleted.
func (c *Client) DeletePrivateKey(i *DeletePrivateKeyInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/private_keys/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}
