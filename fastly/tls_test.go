package fastly

import "testing"

func TestClient_PrivateKey(t *testing.T) {
	t.Parallel()

	fixtureBase := "tls/"

	// Create
	var err error
	var pk *PrivateKey
	record(t, fixtureBase+"create", func(c *Client) {
		pk, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  String("-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n"),
			Name: String("My private key"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeletePrivateKey(&DeletePrivateKeyInput{
				ID: pk.ID,
			})
		})
	}()

	// List
	var lpk []*PrivateKey
	record(t, fixtureBase+"list", func(c *Client) {
		lpk, err = c.ListPrivateKeys(&ListPrivateKeysInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lpk) < 1 {
		t.Errorf("bad privatekeys: %v", lpk)
	}

	// Get
	var gpk *PrivateKey
	record(t, fixtureBase+"get", func(c *Client) {
		gpk, err = c.GetPrivateKey(&GetPrivateKeyInput{
			ID: pk.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if pk.Name != gpk.Name {
		t.Errorf("bad name: %q (%q)", pk.Name, gpk.Name)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: pk.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPrivateKeys_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tls/list", func(c *Client) {
		_, err = c.ListPrivateKeys(&ListPrivateKeysInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetPrivateKey_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tls/get", func(c *Client) {
		_, err = c.GetPrivateKey(&GetPrivateKeyInput{
			ID: "PRIVATE_KEY_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreatePrivateKey_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tls/create", func(c *Client) {
		_, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  String("-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n"),
			Name: String("My private key"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeletePrivateKey_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tls/delete", func(c *Client) {
		err = c.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: "PRIVATE_KEY_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
