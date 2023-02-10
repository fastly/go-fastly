package fastly

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"sort"
	"testing"
)

func TestClient_CreateSecretStore(t *testing.T) {
	t.Parallel()

	var (
		ss  *SecretStore
		err error
	)
	record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		ss, err = c.CreateSecretStore(&CreateSecretStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret store: %v", err)
	}

	// Ensure Secret Store is cleaned up.
	t.Cleanup(func() {
		record(t, fmt.Sprintf("secret_store/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteSecretStore(&DeleteSecretStoreInput{
				ID: ss.ID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting secret store %q: %v", ss.ID, err)
		}
	})

	if got := ss.ID; len(got) == 0 {
		t.Errorf("ID: got %q, want not empty", got)
	}
	if got, want := ss.Name, t.Name(); got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
}

func TestClient_ListSecretStores(t *testing.T) {
	// Cannot be run in parallel, since the list of stores is singular
	// for test user.

	var (
		stores []*SecretStore
		err    error
	)
	for i := 0; i < 5; i++ {
		ss := createSecretStoreHelper(t, i)
		stores = append(stores, ss)
	}
	sort.Slice(stores, func(i, j int) bool {
		return stores[i].ID < stores[j].ID
	})

	var list *SecretStores
	record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		list, err = c.ListSecretStores(&ListSecretStoresInput{})
	})
	if err != nil {
		t.Fatalf("error listing secret store: %v", err)
	}

	if got, want := len(list.Data), len(stores); got != want {
		t.Fatalf("Data: got length %d, want %d", got, want)
	}
	sort.Slice(list.Data, func(i, j int) bool {
		return list.Data[i].ID < list.Data[j].ID
	})
	for i, ss := range list.Data {
		if got, want := ss.ID, stores[i].ID; got != want {
			t.Errorf("Data[%d].ID: got %q, %q", i, got, want)
		}
		if got, want := ss.Name, stores[i].Name; got != want {
			t.Errorf("Data[%d].Name: got %q, %q", i, got, want)
		}
	}

	if got, wantMin := list.Meta.Limit, len(stores); got < wantMin {
		t.Errorf("Meta.Limit: got %d, want >= %d", got, wantMin)
	}
	// Only a single page of results is expected.
	if got, want := list.Meta.NextCursor, ""; got != want {
		t.Errorf("Meta.NextCursor: got %q, want %q", got, want)
	}
}

func TestClient_GetSecretStore(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		store *SecretStore
		err   error
	)
	record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		store, err = c.GetSecretStore(&GetSecretStoreInput{
			ID: ss.ID,
		})
	})
	if err != nil {
		t.Fatalf("error getting secret store: %v", err)
	}

	if got, want := store.ID, ss.ID; got != want {
		t.Errorf("ID: got %q, want %q", got, want)
	}
	if got, want := store.Name, ss.Name; got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
}

func TestClient_DeleteSecretStore(t *testing.T) {
	t.Parallel()

	var (
		ss  *SecretStore
		err error
	)
	record(t, fmt.Sprintf("secret_store/%s/create_store", t.Name()), func(c *Client) {
		ss, err = c.CreateSecretStore(&CreateSecretStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret store: %v", err)
	}

	record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		err = c.DeleteSecretStore(&DeleteSecretStoreInput{
			ID: ss.ID,
		})
	})
	if err != nil {
		t.Fatalf("error deleting secret store: %v", err)
	}
}

func TestClient_CreateSecret(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		s   *Secret
		err error
	)
	record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			ID:     ss.ID,
			Name:   t.Name(),
			Secret: []byte("secretum servare"),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret: %v", err)
	}

	if got, want := s.Name, t.Name(); got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
	if got := s.Digest; len(got) == 0 {
		t.Errorf("Digest: got %q, want not blank", string(got))
	}
}

func TestClient_CreateSecret_clientEncryption(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		ck  *ClientKey
		err error
	)

	record(t, fmt.Sprintf("secret_store/%s/create_client_key", t.Name()), func(c *Client) {
		ck, err = c.CreateClientKey()
	})
	if err != nil {
		t.Fatalf("error creating client key: %v", err)
	}

	if got := ck.PublicKey; len(got) == 0 {
		t.Errorf("PublicKey: got empty")
	}

	if got := ck.Signature; len(got) == 0 {
		t.Errorf("Signature: got empty")
	}

	if got := ck.ExpiresAt; got.IsZero() {
		t.Errorf("ExpiresAt: got empty")
	}

	var sk ed25519.PublicKey

	record(t, fmt.Sprintf("secret_store/%s/get_signing_key", t.Name()), func(c *Client) {
		sk, err = c.GetSigningKey()
	})
	if err != nil {
		t.Fatalf("error getting signing key: %v", err)
	}

	if len(sk) == 0 {
		t.Fatalf("got empty signing key")
	}

	if !ck.VerifySignature(sk) {
		t.Fatalf("signature validation failed")
	}

	enc, err := ck.Encrypt([]byte("secretum servare"))
	if err != nil {
		t.Fatalf("error locally encrypting secret: %v", err)
	}

	var s *Secret

	record(t, fmt.Sprintf("secret_store/%s/create_secret", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			ID:        ss.ID,
			Name:      t.Name(),
			ClientKey: ck.PublicKey,
			Secret:    enc,
		})
	})
	if err != nil {
		t.Fatalf("error creating secret: %v", err)
	}

	if got, want := s.Name, t.Name(); got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
	if got := s.Digest; len(got) == 0 {
		t.Errorf("Digest: got %q, want not blank", string(got))
	}
}

func TestClient_ListSecrets(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		secrets []*Secret
		s       *Secret
		err     error
	)
	for i := 0; i < 5; i++ {
		record(t, fmt.Sprintf("secret_store/%s/create_secret_%02d", t.Name(), i), func(c *Client) {
			s, err = c.CreateSecret(&CreateSecretInput{
				ID:     ss.ID,
				Name:   fmt.Sprintf("%s-%02d", t.Name(), i),
				Secret: []byte("secretum servare"),
			})
		})
		if err != nil {
			t.Fatalf("error creating secret: %v", err)
		}
		secrets = append(secrets, s)
	}

	var list *Secrets
	record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		list, err = c.ListSecrets(&ListSecretsInput{
			ID: ss.ID,
		})
	})
	if err != nil {
		t.Fatalf("error listing secrets: %v", err)
	}

	if got, want := len(list.Data), len(secrets); got != want {
		t.Fatalf("Data: got length %d, want %d", got, want)
	}
	sort.Slice(list.Data, func(i, j int) bool {
		return list.Data[i].Name < list.Data[j].Name
	})
	for i, s := range list.Data {
		if got, want := s.Name, secrets[i].Name; got != want {
			t.Errorf("Data[%d].Name: got %q, %q", i, got, want)
		}
		if got, want := s.Digest, secrets[i].Digest; !bytes.Equal(got, want) {
			t.Errorf("Data[%d].Digest: got %q, %q", i, string(got), string(want))
		}
	}

	if got, wantMin := list.Meta.Limit, len(secrets); got < wantMin {
		t.Errorf("Meta.Limit: got %d, want >= %d", got, wantMin)
	}
	// Only a single page of results is expected.
	if got, want := list.Meta.NextCursor, ""; got != want {
		t.Errorf("Meta.NextCursor: got %q, want %q", got, want)
	}
}

func TestClient_GetSecret(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		s   *Secret
		err error
	)
	record(t, fmt.Sprintf("secret_store/%s/create_secret", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			ID:     ss.ID,
			Name:   t.Name(),
			Secret: []byte("secretum servare"),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret: %v", err)
	}

	var secret *Secret
	record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		secret, err = c.GetSecret(&GetSecretInput{
			ID:   ss.ID,
			Name: s.Name,
		})
	})
	if err != nil {
		t.Fatalf("error getting secret: %v", err)
	}

	if got, want := secret.Name, s.Name; got != want {
		t.Errorf("Name: got %q, %q", got, want)
	}
	if got, want := secret.Digest, s.Digest; !bytes.Equal(got, want) {
		t.Errorf("Digest: got %q, %q", string(got), string(want))
	}
}

func TestClient_DeleteSecret(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		s   *Secret
		err error
	)
	record(t, fmt.Sprintf("secret_store/%s/create_secret", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			ID:     ss.ID,
			Name:   t.Name(),
			Secret: []byte("secretum servare"),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret: %v", err)
	}

	record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		err = c.DeleteSecret(&DeleteSecretInput{
			ID:   ss.ID,
			Name: s.Name,
		})
	})
	if err != nil {
		t.Fatalf("error deleting secret: %v", err)
	}
}

func TestClient_SecretStore_validation(t *testing.T) {
	t.Parallel()

	var err error

	_, err = testClient.CreateSecretStore(&CreateSecretStoreInput{
		Name: "",
	})
	if want := ErrMissingName; err != want {
		t.Errorf("CreateSecretStore: got error %v, want %v", err, want)
	}

	_, err = testClient.GetSecretStore(&GetSecretStoreInput{
		ID: "",
	})
	if want := ErrMissingID; err != want {
		t.Errorf("GetSecretStore: got error %v, want %v", err, want)
	}

	err = testClient.DeleteSecretStore(&DeleteSecretStoreInput{
		ID: "",
	})
	if want := ErrMissingID; err != want {
		t.Errorf("DeleteSecretStore: got error %v, want %v", err, want)
	}

	_, err = testClient.CreateSecret(&CreateSecretInput{
		ID:     "",
		Name:   "name",
		Secret: []byte("secret"),
	})
	if want := ErrMissingID; err != want {
		t.Errorf("CreateSecret: got error %v, want %v", err, want)
	}
	_, err = testClient.CreateSecret(&CreateSecretInput{
		ID:     "123",
		Name:   "",
		Secret: []byte("secret"),
	})
	if want := ErrMissingName; err != want {
		t.Errorf("CreateSecret: got error %v, want %v", err, want)
	}
	_, err = testClient.CreateSecret(&CreateSecretInput{
		ID:     "123",
		Name:   "name",
		Secret: []byte(nil),
	})
	if want := ErrMissingSecret; err != want {
		t.Errorf("CreateSecret: got error %v, want %v", err, want)
	}

	_, err = testClient.ListSecrets(&ListSecretsInput{
		ID: "",
	})
	if want := ErrMissingID; err != want {
		t.Errorf("ListSecrets: got error %v, want %v", err, want)
	}

	_, err = testClient.GetSecret(&GetSecretInput{
		ID:   "",
		Name: "name",
	})
	if want := ErrMissingID; err != want {
		t.Errorf("GetSecret: got error %v, want %v", err, want)
	}
	_, err = testClient.GetSecret(&GetSecretInput{
		ID:   "id",
		Name: "",
	})
	if want := ErrMissingName; err != want {
		t.Errorf("GetSecret: got error %v, want %v", err, want)
	}

	err = testClient.DeleteSecret(&DeleteSecretInput{
		ID:   "",
		Name: "name",
	})
	if want := ErrMissingID; err != want {
		t.Errorf("DeleteSecret: got error %v, want %v", err, want)
	}
	err = testClient.DeleteSecret(&DeleteSecretInput{
		ID:   "id",
		Name: "",
	})
	if want := ErrMissingName; err != want {
		t.Errorf("DeleteSecret: got error %v, want %v", err, want)
	}
}

func createSecretStoreHelper(t *testing.T, i int) *SecretStore {
	t.Helper()

	var (
		ss  *SecretStore
		err error
	)
	record(t, fmt.Sprintf("secret_store/%s/create_store_%02d", t.Name(), i), func(c *Client) {
		ss, err = c.CreateSecretStore(&CreateSecretStoreInput{
			Name: fmt.Sprintf("%s-%02d", t.Name(), i),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret store: %v", err)
	}

	// Cleanup secret store.
	t.Cleanup(func() {
		record(t, fmt.Sprintf("secret_store/%s/delete_store_%02d", t.Name(), i), func(c *Client) {
			err = c.DeleteSecretStore(&DeleteSecretStoreInput{
				ID: ss.ID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting secret store %q: %v", ss.ID, err)
		}
	})

	return ss
}
