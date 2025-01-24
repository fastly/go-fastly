package fastly

import (
	"errors"
	"testing"
)

func TestClient_CustomTLSCertificate(t *testing.T) {
	t.Parallel()

	fixtureBase := "custom_tls/"

	// prepare test key and cert
	privKey, key, err := buildPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	cert, err := buildCertificate(privKey, "example.com")
	if err != nil {
		t.Fatal(err)
	}

	// Create
	var pk *PrivateKey
	Record(t, fixtureBase+"create-key", func(c *Client) {
		pk, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  key,
			Name: "My private key",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create
	var cc *CustomTLSCertificate
	Record(t, fixtureBase+"create", func(c *Client) {
		cc, err = c.CreateCustomTLSCertificate(&CreateCustomTLSCertificateInput{
			CertBlob: cert,
			Name:     "My certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		_ = TestClient.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{
			ID: cc.ID,
		})
		_ = TestClient.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: pk.ID,
		})
	}()

	// List
	var lcc []*CustomTLSCertificate
	Record(t, fixtureBase+"list", func(c *Client) {
		lcc, err = c.ListCustomTLSCertificates(&ListCustomTLSCertificatesInput{
			// NOTE: We set to an explicit false to avoid a test error.
			// This is because we don't activate a real TLS certificate in the test.
			// Filtering by active certs would return zero results from the API call.
			FilterInUse: ToPointer(false),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lcc) < 1 {
		t.Errorf("bad Custom certificates: %v", lcc)
	}

	// Get
	var gcc *CustomTLSCertificate
	Record(t, fixtureBase+"get", func(c *Client) {
		gcc, err = c.GetCustomTLSCertificate(&GetCustomTLSCertificateInput{
			ID: cc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if cc.ID != gcc.ID {
		t.Errorf("bad ID: %q (%q)", cc.ID, gcc.ID)
	}
	if gcc.Domains == nil {
		t.Errorf("Domains should not be nil: %v", cc.Domains)
	}
	if len(gcc.Domains) < 1 {
		t.Errorf("Domains should not be an empty slice: %v", cc.Domains)
	}
	if cc.Domains[0].ID != gcc.Domains[0].ID {
		t.Errorf("bad Domain ID: %q (%q)", cc.Domains[0].ID, gcc.Domains[0].ID)
	}

	// regenerate test cert using the created key above
	cert, err = buildCertificate(privKey, "example.com", "foo.example.com")
	if err != nil {
		t.Fatal(err)
	}

	// Update
	var ucc *CustomTLSCertificate
	Record(t, fixtureBase+"update", func(c *Client) {
		ucc, err = c.UpdateCustomTLSCertificate(&UpdateCustomTLSCertificateInput{
			ID:       cc.ID,
			CertBlob: cert,
			Name:     "My certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if cc.ID != ucc.ID {
		t.Errorf("bad ID: %q (%q)", cc.ID, ucc.ID)
	}

	// Delete
	Record(t, fixtureBase+"delete-cert", func(c *Client) {
		err = c.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{
			ID: cc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	Record(t, fixtureBase+"delete-key", func(c *Client) {
		err = c.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: pk.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateCustomTLSCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	_, err = TestClient.CreateCustomTLSCertificate(&CreateCustomTLSCertificateInput{
		Name: "My certificate",
	})
	if !errors.Is(err, ErrMissingCertBlob) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCustomTLSCertificate_validation(t *testing.T) {
	t.Parallel()

	err := TestClient.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ListCustomTLSCertificates_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "custom_tls/list", func(c *Client) {
		_, err = c.ListCustomTLSCertificates(&ListCustomTLSCertificatesInput{
			FilterInUse: ToPointer(false),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetCustomTLSCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	_, err = TestClient.GetCustomTLSCertificate(&GetCustomTLSCertificateInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCustomTLSCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	_, err = TestClient.UpdateCustomTLSCertificate(&UpdateCustomTLSCertificateInput{
		CertBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
		Name:     "My certificate",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateCustomTLSCertificate(&UpdateCustomTLSCertificateInput{
		ID:   "CERTIFICATE_ID",
		Name: "My certificate",
	})
	if !errors.Is(err, ErrMissingCertBlob) {
		t.Errorf("bad error: %s", err)
	}
}
