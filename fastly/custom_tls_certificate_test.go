package fastly

import (
	"testing"
)

func TestClient_CustomTLSCertificate(t *testing.T) {
	t.Parallel()

	fixtureBase := "custom_tls/"

	// Create
	var err error
	var cc *CustomTLSCertificate
	record(t, fixtureBase+"create", func(c *Client) {
		cc, err = c.CreateCustomTLSCertificate(&CreateCustomTLSCertificateInput{
			CertBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			Name:     "My certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{
				ID: cc.ID,
			})
		})
	}()

	// List
	var lcc []*CustomTLSCertificate
	record(t, fixtureBase+"list", func(c *Client) {
		lcc, err = c.ListCustomTLSCertificates(&ListCustomTLSCertificatesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lcc) < 1 {
		t.Errorf("bad Custom certificates: %v", lcc)
	}

	// Get
	var gcc *CustomTLSCertificate
	record(t, fixtureBase+"get", func(c *Client) {
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

	// Update
	var ucc *CustomTLSCertificate
	record(t, fixtureBase+"update", func(c *Client) {
		ucc, err = c.UpdateCustomTLSCertificate(&UpdateCustomTLSCertificateInput{
			ID:       "CERTIFICATE_ID",
			CertBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
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
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{
			ID: cc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateCustomTLSCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/create", func(c *Client) {
		_, err = c.CreateCustomTLSCertificate(&CreateCustomTLSCertificateInput{
			CertBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			Name:     "My certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.CreateCustomTLSCertificate(&CreateCustomTLSCertificateInput{
		Name: "My certificate",
	})
	if err != ErrMissingCertBlob {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCustomTLSCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/delete", func(c *Client) {
		err = c.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{
			ID: "CERTIFICATE_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	err = testClient.DeleteCustomTLSCertificate(&DeleteCustomTLSCertificateInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ListCustomTLSCertificates_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/list", func(c *Client) {
		_, err = c.ListCustomTLSCertificates(&ListCustomTLSCertificatesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetCustomTLSCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/get", func(c *Client) {
		_, err = c.GetCustomTLSCertificate(&GetCustomTLSCertificateInput{
			ID: "CERTIFICATE_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.GetCustomTLSCertificate(&GetCustomTLSCertificateInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCustomTLSCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/update", func(c *Client) {
		_, err = c.UpdateCustomTLSCertificate(&UpdateCustomTLSCertificateInput{
			ID:       "CERTIFICATE_ID",
			CertBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			Name:     "My certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.UpdateCustomTLSCertificate(&UpdateCustomTLSCertificateInput{
		CertBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
		Name:     "My certificate",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCustomTLSCertificate(&UpdateCustomTLSCertificateInput{
		ID:   "CERTIFICATE_ID",
		Name: "My certificate",
	})
	if err != ErrMissingCertBlob {
		t.Errorf("bad error: %s", err)
	}
}
