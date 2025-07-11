package fastly

import (
	"context"
	"errors"
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
	Record(t, fixtureBase+"create-key", func(c *Client) {
		pk, err = c.CreatePrivateKey(context.TODO(), &CreatePrivateKeyInput{
			Key:  key,
			Name: "My private key",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a customer TLS certificate to pass to the mutual authentication endpoint.
	var cc *CustomTLSCertificate
	Record(t, fixtureBase+"create-cert", func(c *Client) {
		cc, err = c.CreateCustomTLSCertificate(context.TODO(), &CreateCustomTLSCertificateInput{
			CertBlob: cert,
			Name:     "My custom certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		_ = TestClient.DeleteCustomTLSCertificate(context.TODO(), &DeleteCustomTLSCertificateInput{
			ID: cc.ID,
		})
		_ = TestClient.DeletePrivateKey(context.TODO(), &DeletePrivateKeyInput{
			ID: pk.ID,
		})
	}()

	// Create mutual authentication using the custom TLS certificate above.
	var tma *TLSMutualAuthentication
	Record(t, fixtureBase+"create-tma", func(c *Client) {
		tma, err = c.CreateTLSMutualAuthentication(context.TODO(), &CreateTLSMutualAuthenticationInput{
			CertBundle: cert,
			Enforced:   false,
			Name:       "My mutual authentication",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		_ = TestClient.DeleteTLSMutualAuthentication(context.TODO(), &DeleteTLSMutualAuthenticationInput{
			ID: tma.ID,
		})
	}()

	if tma.Enforced {
		t.Errorf("bad Enforced: %t", tma.Enforced)
	}

	// List
	var tmas []*TLSMutualAuthentication
	Record(t, fixtureBase+"list", func(c *Client) {
		tmas, err = c.ListTLSMutualAuthentication(context.TODO(), &ListTLSMutualAuthenticationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tmas) < 1 {
		t.Errorf("failed to return any objects: %v", tmas)
	}

	// Get
	var gtma *TLSMutualAuthentication
	Record(t, fixtureBase+"get", func(c *Client) {
		gtma, err = c.GetTLSMutualAuthentication(context.TODO(), &GetTLSMutualAuthenticationInput{
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
	Record(t, fixtureBase+"update", func(c *Client) {
		utma, err = c.UpdateTLSMutualAuthentication(context.TODO(), &UpdateTLSMutualAuthenticationInput{
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
	Record(t, fixtureBase+"delete-tma", func(c *Client) {
		err = c.DeleteTLSMutualAuthentication(context.TODO(), &DeleteTLSMutualAuthenticationInput{
			ID: tma.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	Record(t, fixtureBase+"delete-cert", func(c *Client) {
		err = c.DeleteCustomTLSCertificate(context.TODO(), &DeleteCustomTLSCertificateInput{
			ID: cc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	Record(t, fixtureBase+"delete-key", func(c *Client) {
		err = c.DeletePrivateKey(context.TODO(), &DeletePrivateKeyInput{
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
	_, err = TestClient.CreateTLSMutualAuthentication(context.TODO(), &CreateTLSMutualAuthenticationInput{
		Name: "My certificate",
	})
	if !errors.Is(err, ErrMissingCertBundle) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	err := TestClient.DeleteTLSMutualAuthentication(context.TODO(), &DeleteTLSMutualAuthenticationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ListTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "mutual_authentication/list", func(c *Client) {
		_, err = c.ListTLSMutualAuthentication(context.TODO(), &ListTLSMutualAuthenticationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	var err error
	_, err = TestClient.GetCustomTLSCertificate(context.TODO(), &GetCustomTLSCertificateInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateTLSMutualAuthentication_validation(t *testing.T) {
	t.Parallel()

	var err error
	_, err = TestClient.UpdateTLSMutualAuthentication(context.TODO(), &UpdateTLSMutualAuthenticationInput{
		ID:   "example",
		Name: "My certificate",
	})
	if !errors.Is(err, ErrMissingCertBundle) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateTLSMutualAuthentication(context.TODO(), &UpdateTLSMutualAuthenticationInput{
		CertBundle: "example",
		Name:       "My certificate",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}
