package fastly

import (
	"testing"
)

func TestClient_TLSMutualAuthentication(t *testing.T) {
	t.Parallel()

	fixtureBase := "mutual_authentication/"

	privateKey, key, err := buildPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	cert, err := buildCertificate(privateKey, "example.com")
	if err != nil {
		t.Fatal(err)
	}

	// Create private key required to generate a custom certificate.
	var pk *PrivateKey
	record(t, fixtureBase+"create-key", func(c *Client) {
		pk, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  key,
			Name: "My private key",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a customer TLS certificate to pass to the mutual authentication endpoint.
	var cc *CustomTLSCertificate
	record(t, fixtureBase+"create-cert", func(c *Client) {
		cc, err = c.CreateCustomTLSCertificate(&CreateCustomTLSCertificateInput{
			CertBlob: cert,
			Name:     "My custom certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		_ = testClient.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{
			ID: cc.ID,
		})
		_ = testClient.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: pk.ID,
		})
	}()

	// Create mutual authentication using the custom TLS certificate above.
	var tma *TLSMutualAuthentication
	record(t, fixtureBase+"create-tma", func(c *Client) {
		tma, err = c.CreateTLSMutualAuthentication(&CreateTLSMutualAuthenticationInput{
			CertBundle: cert,
			Name:       "My mutual authentication",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		_ = testClient.DeleteTLSMutualAuthentication(&DeleteTLSMutualAuthenticationInput{
			ID: tma.ID,
		})
	}()

	// List
	var tmas []*TLSMutualAuthentication
	record(t, fixtureBase+"list", func(c *Client) {
		tmas, err = c.ListTLSMutualAuthentication(&ListTLSMutualAuthenticationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tmas) < 1 {
		t.Errorf("failed to return any objects: %v", tmas)
	}

	// Get
	var gtma *TLSMutualAuthentication
	record(t, fixtureBase+"get", func(c *Client) {
		gtma, err = c.GetTLSMutualAuthentication(&GetTLSMutualAuthenticationInput{
			ID: tma.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if tma.ID != gtma.ID {
		t.Errorf("bad ID: %q (%q)", tma.ID, gtma.ID)
	}

	// Update
	var utma *TLSMutualAuthentication
	record(t, fixtureBase+"update", func(c *Client) {
		utma, err = c.UpdateTLSMutualAuthentication(&UpdateTLSMutualAuthenticationInput{
			CertBundle: cert,
			Enforced:   true,
			ID:         tma.ID,
			Name:       "My mutual authentication updated",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if tma.ID != utma.ID {
		t.Errorf("bad ID: %q (%q)", tma.ID, utma.ID)
	}

	// Delete
	record(t, fixtureBase+"delete-tma", func(c *Client) {
		err = c.DeleteTLSMutualAuthentication(&DeleteTLSMutualAuthenticationInput{
			ID: tma.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	record(t, fixtureBase+"delete-cert", func(c *Client) {
		err = c.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{
			ID: cc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	record(t, fixtureBase+"delete-key", func(c *Client) {
		err = c.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: pk.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	var err error
	_, err = testClient.CreateTLSMutualAuthentication(&CreateTLSMutualAuthenticationInput{
		Name: "My certificate",
	})
	if err != ErrMissingCertBundle {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	err := testClient.DeleteTLSMutualAuthentication(&DeleteTLSMutualAuthenticationInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ListTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "mutual_authentication/list", func(c *Client) {
		_, err = c.ListTLSMutualAuthentication(&ListTLSMutualAuthenticationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	var err error
	_, err = testClient.GetCustomTLSCertificate(&GetCustomTLSCertificateInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	var err error
	_, err = testClient.UpdateTLSMutualAuthentication(&UpdateTLSMutualAuthenticationInput{
		ID:   "example",
		Name: "My certificate",
	})
	if err != ErrMissingCertBundle {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateTLSMutualAuthentication(&UpdateTLSMutualAuthenticationInput{
		CertBundle: "example",
		Name:       "My certificate",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
