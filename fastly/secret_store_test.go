package fastly

import (
	"bytes"
	"crypto/ed25519"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"testing"
)

func TestClient_CreateSecretStore(t *testing.T) {
	t.Parallel()

	var (
		ss  *SecretStore
		err error
	)
	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		ss, err = c.CreateSecretStore(&CreateSecretStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret store: %v", err)
	}

	// Ensure Secret Store is cleaned up.
	t.Cleanup(func() {
		Record(t, fmt.Sprintf("secret_store/%s/delete_store", t.Name()), func(c *Client) {
			err = c.DeleteSecretStore(&DeleteSecretStoreInput{
				StoreID: ss.StoreID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting secret store %q: %v", ss.StoreID, err)
		}
	})

	if got := ss.StoreID; len(got) == 0 {
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
		return stores[i].StoreID < stores[j].StoreID
	})

	var list *SecretStores
	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		list, err = c.ListSecretStores(&ListSecretStoresInput{})
	})
	if err != nil {
		t.Fatalf("error listing secret store: %v", err)
	}

	if got, want := len(list.Data), len(stores); got != want {
		t.Fatalf("Data: got length %d, want %d", got, want)
	}
	sort.Slice(list.Data, func(i, j int) bool {
		return list.Data[i].StoreID < list.Data[j].StoreID
	})
	for i, ss := range list.Data {
		if got, want := ss.StoreID, stores[i].StoreID; got != want {
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

	Record(t, fmt.Sprintf("secret_store/%s/list-with-name", t.Name()), func(c *Client) {
		list, err = c.ListSecretStores(&ListSecretStoresInput{Name: stores[0].Name})
	})

	if err != nil {
		t.Fatalf("error listing secret store by name: %v", err)
	}

	if got, want := len(list.Data), 1; got != want {
		t.Fatalf("Data: got length %d, want %d", got, want)
	}

	if got, want := list.Data[0].StoreID, stores[0].StoreID; got != want {
		t.Errorf("Data[0].ID: got %q, want %q", got, want)
	}
}

func TestClient_GetSecretStore(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		store *SecretStore
		err   error
	)
	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		store, err = c.GetSecretStore(&GetSecretStoreInput{
			StoreID: ss.StoreID,
		})
	})
	if err != nil {
		t.Fatalf("error getting secret store: %v", err)
	}

	if got, want := store.StoreID, ss.StoreID; got != want {
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
	Record(t, fmt.Sprintf("secret_store/%s/create_store", t.Name()), func(c *Client) {
		ss, err = c.CreateSecretStore(&CreateSecretStoreInput{
			Name: t.Name(),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret store: %v", err)
	}

	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		err = c.DeleteSecretStore(&DeleteSecretStoreInput{
			StoreID: ss.StoreID,
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
	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			StoreID: ss.StoreID,
			Name:    t.Name(),
			Secret:  []byte("secretum servare"),
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

func TestClient_CreateOrRecreateSecret(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		s   *Secret
		err error
	)
	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			StoreID: ss.StoreID,
			Name:    t.Name(),
			Secret:  []byte("secretum servare"),
			Method:  http.MethodPut,
		})
	})
	if err != nil {
		t.Fatalf("error creating or recreating secret: %v", err)
	}

	if got, want := s.Name, t.Name(); got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
	if got := s.Digest; len(got) == 0 {
		t.Errorf("Digest: got %q, want not blank", string(got))
	}
	if got, want := s.Recreated, false; got != want {
		t.Errorf("Recreated: got %v, want %v", got, want)
	}
}

func TestClient_RecreateSecret(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		s   *Secret
		err error
	)
	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		// There must be an existing secret already, otherwise
		// the following PATCH request will fail.
		s, err = c.CreateSecret(&CreateSecretInput{
			StoreID: ss.StoreID,
			Name:    t.Name(),
			Secret:  []byte("secretum servare"),
		})
		if err != nil {
			return
		}

		s, err = c.CreateSecret(&CreateSecretInput{
			StoreID: ss.StoreID,
			Name:    t.Name(),
			Secret:  []byte("secretum servare"),
			Method:  http.MethodPatch,
		})
	})
	if err != nil {
		t.Fatalf("error recreating secret: %v", err)
	}

	if got, want := s.Name, t.Name(); got != want {
		t.Errorf("Name: got %q, want %q", got, want)
	}
	if got := s.Digest; len(got) == 0 {
		t.Errorf("Digest: got %q, want not blank", string(got))
	}
	if got, want := s.Recreated, true; got != want {
		t.Errorf("Recreated: got %v, want %v", got, want)
	}
}

func TestClient_CreateSecret_clientEncryption(t *testing.T) {
	t.Parallel()

	ss := createSecretStoreHelper(t, 0)

	var (
		ck  *ClientKey
		err error
	)

	Record(t, fmt.Sprintf("secret_store/%s/create_client_key", t.Name()), func(c *Client) {
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

	Record(t, fmt.Sprintf("secret_store/%s/get_signing_key", t.Name()), func(c *Client) {
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

	Record(t, fmt.Sprintf("secret_store/%s/create_secret", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			StoreID:   ss.StoreID,
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
		Record(t, fmt.Sprintf("secret_store/%s/create_secret_%02d", t.Name(), i), func(c *Client) {
			s, err = c.CreateSecret(&CreateSecretInput{
				StoreID: ss.StoreID,
				Name:    fmt.Sprintf("%s-%02d", t.Name(), i),
				Secret:  []byte("secretum servare"),
			})
		})
		if err != nil {
			t.Fatalf("error creating secret: %v", err)
		}
		secrets = append(secrets, s)
	}

	var list *Secrets
	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		list, err = c.ListSecrets(&ListSecretsInput{
			StoreID: ss.StoreID,
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
	Record(t, fmt.Sprintf("secret_store/%s/create_secret", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			StoreID: ss.StoreID,
			Name:    t.Name(),
			Secret:  []byte("secretum servare"),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret: %v", err)
	}

	var secret *Secret
	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		secret, err = c.GetSecret(&GetSecretInput{
			StoreID: ss.StoreID,
			Name:    s.Name,
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
	Record(t, fmt.Sprintf("secret_store/%s/create_secret", t.Name()), func(c *Client) {
		s, err = c.CreateSecret(&CreateSecretInput{
			StoreID: ss.StoreID,
			Name:    t.Name(),
			Secret:  []byte("secretum servare"),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret: %v", err)
	}

	Record(t, fmt.Sprintf("secret_store/%s", t.Name()), func(c *Client) {
		err = c.DeleteSecret(&DeleteSecretInput{
			StoreID: ss.StoreID,
			Name:    s.Name,
		})
	})
	if err != nil {
		t.Fatalf("error deleting secret: %v", err)
	}
}

func TestClient_SecretStore_validation(t *testing.T) {
	t.Parallel()

	var err error

	_, err = TestClient.CreateSecretStore(&CreateSecretStoreInput{
		Name: "",
	})
	if want := ErrMissingName; !errors.Is(want, err) {
		t.Errorf("CreateSecretStore: got error %v, want %v", err, want)
	}

	_, err = TestClient.GetSecretStore(&GetSecretStoreInput{
		StoreID: "",
	})
	if want := ErrMissingStoreID; !errors.Is(want, err) {
		t.Errorf("GetSecretStore: got error %v, want %v", err, want)
	}

	err = TestClient.DeleteSecretStore(&DeleteSecretStoreInput{
		StoreID: "",
	})
	if want := ErrMissingStoreID; !errors.Is(want, err) {
		t.Errorf("DeleteSecretStore: got error %v, want %v", err, want)
	}

	_, err = TestClient.CreateSecret(&CreateSecretInput{
		StoreID: "",
		Name:    "name",
		Secret:  []byte("secret"),
	})
	if want := ErrMissingStoreID; !errors.Is(want, err) {
		t.Errorf("CreateSecret: got error %v, want %v", err, want)
	}
	_, err = TestClient.CreateSecret(&CreateSecretInput{
		StoreID: "123",
		Name:    "",
		Secret:  []byte("secret"),
	})
	if want := ErrMissingName; !errors.Is(want, err) {
		t.Errorf("CreateSecret: got error %v, want %v", err, want)
	}
	_, err = TestClient.CreateSecret(&CreateSecretInput{
		StoreID: "123",
		Name:    "name",
		Secret:  []byte(nil),
	})
	if want := ErrMissingSecret; !errors.Is(want, err) {
		t.Errorf("CreateSecret: got error %v, want %v", err, want)
	}

	_, err = TestClient.ListSecrets(&ListSecretsInput{
		StoreID: "",
	})
	if want := ErrMissingStoreID; !errors.Is(want, err) {
		t.Errorf("ListSecrets: got error %v, want %v", err, want)
	}

	_, err = TestClient.GetSecret(&GetSecretInput{
		StoreID: "",
		Name:    "name",
	})
	if want := ErrMissingStoreID; !errors.Is(want, err) {
		t.Errorf("GetSecret: got error %v, want %v", err, want)
	}
	_, err = TestClient.GetSecret(&GetSecretInput{
		StoreID: "id",
		Name:    "",
	})
	if want := ErrMissingName; !errors.Is(want, err) {
		t.Errorf("GetSecret: got error %v, want %v", err, want)
	}

	err = TestClient.DeleteSecret(&DeleteSecretInput{
		StoreID: "",
		Name:    "name",
	})
	if want := ErrMissingStoreID; !errors.Is(want, err) {
		t.Errorf("DeleteSecret: got error %v, want %v", err, want)
	}
	err = TestClient.DeleteSecret(&DeleteSecretInput{
		StoreID: "id",
		Name:    "",
	})
	if want := ErrMissingName; !errors.Is(want, err) {
		t.Errorf("DeleteSecret: got error %v, want %v", err, want)
	}
}

func createSecretStoreHelper(t *testing.T, i int) *SecretStore {
	t.Helper()

	var (
		ss  *SecretStore
		err error
	)
	Record(t, fmt.Sprintf("secret_store/%s/create_store_%02d", t.Name(), i), func(c *Client) {
		ss, err = c.CreateSecretStore(&CreateSecretStoreInput{
			Name: fmt.Sprintf("%s-%02d", t.Name(), i),
		})
	})
	if err != nil {
		t.Fatalf("error creating secret store: %v", err)
	}

	// Cleanup secret store.
	t.Cleanup(func() {
		Record(t, fmt.Sprintf("secret_store/%s/delete_store_%02d", t.Name(), i), func(c *Client) {
			err = c.DeleteSecretStore(&DeleteSecretStoreInput{
				StoreID: ss.StoreID,
			})
		})
		if err != nil {
			t.Fatalf("error deleting secret store %q: %v", ss.StoreID, err)
		}
	})

	return ss
}
