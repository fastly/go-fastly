package fastly

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/nacl/box"
)

// Secret Store.
// A secret store is a persistent, globally distributed store for secrets accessible to Compute@Edge services during request processing.
// https://developer.fastly.com/reference/api/secret-store/

// SecretStoreMeta is the metadata returned from Secret Store paginated responses.
type SecretStoreMeta struct {
	// Limit is the limit of results returned.
	Limit int `json:"limit"`
	// NextCursor can be used on a subsequent request to fetch the next page of data.
	NextCursor string `json:"next_cursor,omitempty"`
}

// SecretStore represents a Secret Store response from the Fastly API.
type SecretStore struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	StoreID   string    `json:"id"`
}

// CreateSecretStoreInput is used as input to the CreateSecretStore function.
type CreateSecretStoreInput struct {
	// Name of the Secret Store (required).
	Name string `json:"name"`
}

// CreateSecretStore creates a new resource.
func (c *Client) CreateSecretStore(i *CreateSecretStoreInput) (*SecretStore, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}

	p := "/resources/stores/secret"
	resp, err := c.PostJSON(p, i, &RequestOptions{
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
	// Cursor is the pagination cursor (optional).
	Cursor string
	// Limit is the desired number of Secret Stores (optional).
	Limit int
	// Name is the name of the secret store (optional).
	Name string
}

// ListSecretStores retrieves all resources.
//
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
	if i.Name != "" {
		params["name"] = i.Name
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
	// StoreID of the Secret Store (required).
	StoreID string
}

// GetSecretStore retrieves the specified resource.
func (c *Client) GetSecretStore(i *GetSecretStoreInput) (*SecretStore, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	p := "/resources/stores/secret/" + i.StoreID

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
	// StoreID of the Secret Store (required).
	StoreID string
}

// DeleteSecretStore deletes the specified resource.
func (c *Client) DeleteSecretStore(i *DeleteSecretStoreInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}

	p := "/resources/stores/secret/" + i.StoreID

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
	Name      string    `json:"name"`
	Digest    []byte    `json:"digest"` // Digest is an opaque hash of the secret.
	CreatedAt time.Time `json:"created_at"`
	Recreated bool      `json:"recreated,omitempty"`
}

// CreateSecretInput is used as input to the CreateSecret function.
type CreateSecretInput struct {
	// ClientKey is the public key used to encrypt the secret with (optional).
	ClientKey []byte
	// Method is the HTTP request method used to create the secret.
	//
	// Secret names must be unique within a store.
	// The method effects how duplicate names are handled:
	//
	// - POST:  Default. Create a secret and error if one already exists with the same name.
	// - PUT:   Create or recreate a secret.
	// - PATCH: Recreate a secret and error if one does not already exist with the same name.
	//
	// More details: https://developer.fastly.com/reference/api/services/resources/secret-store-secret/
	Method string
	// Name of the Secret (required).
	Name string
	// Secret is the plaintext secret to be stored (required).
	// The value will be base64-encoded when delivered to the API, which is the
	// required format.
	Secret []byte
	// StoreID of the Secret Store (required).
	StoreID string
}

// CreateSecret creates a new resource.
func (c *Client) CreateSecret(i *CreateSecretInput) (*Secret, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if len(i.Secret) == 0 {
		return nil, ErrMissingSecret
	}

	p := "/resources/stores/secret/" + i.StoreID + "/secrets"

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(struct {
		Name      string `json:"name"`
		Secret    []byte `json:"secret"`
		ClientKey []byte `json:"client_key,omitempty"`
	}{
		Name:      i.Name,
		Secret:    i.Secret,
		ClientKey: i.ClientKey,
	})
	if err != nil {
		return nil, err
	}

	method := i.Method
	if method == "" {
		method = http.MethodPost
	}
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		// Method is allowed.
	default:
		return nil, ErrInvalidMethod
	}

	resp, err := c.Request(method, p, &RequestOptions{
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
	// Cursor is the pagination cursor (optional).
	Cursor string
	// Limit is the desired number of Secrets (optional).
	Limit int
	// StoreID of the Secret Store (required).
	StoreID string
}

// ListSecrets retrieves all resources.
//
// The returned next cursor, if non-blank, can be used as input to a subsequent
// request for the next page of results.
func (c *Client) ListSecrets(i *ListSecretsInput) (*Secrets, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}

	p := "/resources/stores/secret/" + i.StoreID + "/secrets"

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
	// Name of the Secret (required).
	Name string
	// StoreID of the Secret Store (required).
	StoreID string
}

// GetSecret retrieves the specified resource.
func (c *Client) GetSecret(i *GetSecretInput) (*Secret, error) {
	if i.StoreID == "" {
		return nil, ErrMissingStoreID
	}
	if i.Name == "" {
		return nil, ErrMissingName
	}

	p := "/resources/stores/secret/" + i.StoreID + "/secrets/" + i.Name

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
	// Name of the secret (required).
	Name string
	// StoreID of the Secret Store (required).
	StoreID string
}

// DeleteSecret deletes the specified resource.
func (c *Client) DeleteSecret(i *DeleteSecretInput) error {
	if i.StoreID == "" {
		return ErrMissingStoreID
	}
	if i.Name == "" {
		return ErrMissingName
	}

	p := "/resources/stores/secret/" + i.StoreID + "/secrets/" + i.Name

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

// ClientKey is an X25519 public key that can be used with
// golang.org/x/crypto/nacl/box to encrypt secrets locally before
// uploading them to the Fastly API.  A client key is valid only for a
// short amount of time, and should be used immediately.  The key is not
// valid after the ExpiresAt time.
//
// Client keys are signed, and the attached signature must be validated
// using the public signing key before it is used.  A ValidateSignature
// method is provided for this purpose.
type ClientKey struct {
	PublicKey []byte    `json:"public_key"`
	Signature []byte    `json:"signature"`
	ExpiresAt time.Time `json:"expires_at"`
}

// VerifySignature reports if the signingKey was used to generate the
// client key's signature.  It must be a valid Ed25519 public key, and
// it will panic if len(signingKey) is not ed25519.PublicKeySize.
// https://pkg.go.dev/crypto/ed25519#PublicKeySize
func (ck *ClientKey) VerifySignature(signingKey ed25519.PublicKey) bool {
	return ed25519.Verify(signingKey, ck.PublicKey, ck.Signature)
}

// Encrypt uses the client key to encrypt the provided plaintext
// using a libsodium-compatible sealed box.
// https://pkg.go.dev/golang.org/x/crypto/nacl/box#SealAnonymous
// https://libsodium.gitbook.io/doc/public-key_cryptography/sealed_boxes
func (ck *ClientKey) Encrypt(plaintext []byte) ([]byte, error) {
	if len(ck.PublicKey) != 32 {
		return nil, fmt.Errorf("invalid public key length %d", len(ck.PublicKey))
	}

	return box.SealAnonymous(nil, plaintext, (*[32]byte)(ck.PublicKey), rand.Reader)
}

// CreateClientKey creates a new time-limited client key for locally
// encrypting secrets before uploading them to the Fastly API.
func (c *Client) CreateClientKey() (*ClientKey, error) {
	p := "/resources/stores/secret/client-key"

	resp, err := c.Post(p, &RequestOptions{
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

	var output ClientKey
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil
}

// GetSigningKey returns the public signing key for client keys.  In
// general the signing key changes very rarely, and it's recommended to
// ship the signing key out-of-band from the API.
func (c *Client) GetSigningKey() (ed25519.PublicKey, error) {
	p := "/resources/stores/secret/signing-key"

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

	var output struct {
		SigningKey []byte `json:"signing_key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return ed25519.PublicKey(output.SigningKey), nil
}
