package fastly

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Secret Store.
// A secret store is a persistent, globally distributed store for secrets accessible to Compute@Edge services during request processing.
// https://developer.fastly.com/reference/api/secret-store/

// SecretStoreMeta is the metadata returned from Secret Store paginated responses.
type SecretStoreMeta struct {
	// Limit is the limit of results returned.
	Limit int `json:"limit"`
	// NextCursor can be used on a subsequent request to fetch
	// the next page of data.
	NextCursor string `json:"next_cursor,omitempty"`
}

// SecretStore represents a Secret Store response from the Fastly API.
type SecretStore struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateSecretStoreInput is used as input to the CreateSecretStore function.
type CreateSecretStoreInput struct {
	// Name of the Secret Store (required).
	Name string
}

// CreateSecretStore creates a new Secret Store.
func (c *Client) CreateSecretStore(i *CreateSecretStoreInput) (*SecretStore, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}

	p := "/resources/stores/secret"

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(struct {
		Name string `json:"name"`
	}{
		Name: i.Name,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(p, &RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Body:     &body,
		Parallel: true,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output SecretStore
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil
}

// SecretStores represents a list of Secret Stores from the Fastly API.
type SecretStores struct {
	// Data is a list of Secret Stores.
	Data []SecretStore `json:"data"`
	// Meta contains response pagination data.
	Meta SecretStoreMeta `json:"meta"`
}

// ListSecretStoresInput is used as input to the ListSecretStores function.
type ListSecretStoresInput struct {
	// Limit is the desired number of Secret Stores (optional).
	Limit int
	// Cursor is the pagination cursor (optional).
	Cursor string
}

// ListSecretStores returns a paginated list of Secret Stores.
// The returned next cursor, if non-blank, can be used as input to a subsequent
// request for the next page of results.
func (c *Client) ListSecretStores(i *ListSecretStoresInput) (*SecretStores, error) {
	p := "/resources/stores/secret"

	params := make(map[string]string, 2)
	if i.Limit > 0 {
		params["limit"] = strconv.Itoa(i.Limit)
	}
	if i.Cursor != "" {
		params["cursor"] = i.Cursor
	}

	resp, err := c.Get(p, &RequestOptions{
		Params: params,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Parallel: true,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output SecretStores
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil
}

// GetSecretStoreInput is used as input to the GetSecretStore function.
type GetSecretStoreInput struct {
	// ID of the Secret Store (required).
	ID string
}

// GetSecretStore gets a single Secret Store.
func (c *Client) GetSecretStore(i *GetSecretStoreInput) (*SecretStore, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := "/resources/stores/secret/" + i.ID

	resp, err := c.Get(p, &RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Parallel: true,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output SecretStore
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil
}

// DeleteSecretStoreInput is used as input to the DeleteSecretStore function.
type DeleteSecretStoreInput struct {
	// ID of the Secret Store (required).
	ID string
}

// DeleteSecretStore deletes the given Secret Store and associated Secrets.
func (c *Client) DeleteSecretStore(i *DeleteSecretStoreInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	p := "/resources/stores/secret/" + i.ID

	resp, err := c.Delete(p, &RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Parallel: true,
	})
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// Secret is a Secret Store secret.
type Secret struct {
	Name string `json:"name"`
	// Digest is an opaque hash of the secret.
	Digest []byte `json:"digest"`
}

// CreateSecretInput is used as input to the CreateSecret function.
type CreateSecretInput struct {
	// ID of the Secret Store (required).
	ID string
	// Name of the Secret (required).
	Name string
	// Secret is the plaintext secret to be stored (required).
	// The value will be base64-encoded when delivered to the API,
	// which is the required format.
	Secret []byte
}

// CreateSecret creates a new Secret within a Secret Store.
func (c *Client) CreateSecret(i *CreateSecretInput) (*Secret, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if len(i.Secret) == 0 {
		return nil, ErrMissingSecret
	}

	p := "/resources/stores/secret/" + i.ID + "/secrets"

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(struct {
		Name   string `json:"name"`
		Secret []byte `json:"secret"`
	}{
		Name:   i.Name,
		Secret: i.Secret,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(p, &RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Body:     &body,
		Parallel: true,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output Secret
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil
}

// Secrets represents a list of Secrets from the Fastly API.
type Secrets struct {
	// Data is a list of Secrets.
	Data []Secret `json:"data"`
	// Meta contains pagination data.
	Meta SecretStoreMeta `json:"meta"`
}

// ListSecretsInput is used as input to the ListSecrets function.
type ListSecretsInput struct {
	// ID of the Secret Store (required).
	ID string
	// Limit is the desired number of Secrets (optional).
	Limit int
	// Cursor is the pagination cursor (optional).
	Cursor string
}

// ListSecrets returns a list of Secrets for the given Secret Store.
// The returned next cursor, if non-blank, can be used as input to a subsequent
// request for the next page of results.
func (c *Client) ListSecrets(i *ListSecretsInput) (*Secrets, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	p := "/resources/stores/secret/" + i.ID + "/secrets"

	params := make(map[string]string, 2)
	if i.Limit > 0 {
		params["limit"] = strconv.Itoa(i.Limit)
	}
	if i.Cursor != "" {
		params["cursor"] = i.Cursor
	}

	resp, err := c.Get(p, &RequestOptions{
		Params: params,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Parallel: true,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output Secrets
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil
}

// GetSecretInput is used as input to the GetSecret function.
type GetSecretInput struct {
	// ID of the Secret Store (required).
	ID string
	// Name of the Secret (required).
	Name string
}

// GetSecret returns a single Secret from a given Secret Store.
func (c *Client) GetSecret(i *GetSecretInput) (*Secret, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}
	if i.Name == "" {
		return nil, ErrMissingName
	}

	p := "/resources/stores/secret/" + i.ID + "/secrets/" + i.Name

	resp, err := c.Get(p, &RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Parallel: true,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output Secret
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil
}

// DeleteSecretInput is used as input to the DeleteSecret function.
type DeleteSecretInput struct {
	// ID of the Secret Store (required).
	ID string
	// Name of the secret (required).
	Name string
}

// DeleteSecret deletes the given Secret.
func (c *Client) DeleteSecret(i *DeleteSecretInput) error {
	if i.ID == "" {
		return ErrMissingID
	}
	if i.Name == "" {
		return ErrMissingName
	}

	p := "/resources/stores/secret/" + i.ID + "/secrets/" + i.Name

	resp, err := c.Delete(p, &RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Parallel: true,
	})
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
