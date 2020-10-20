package fastly

import "testing"

func TestClient_CustomCertificate(t *testing.T) {
	t.Parallel()

	fixtureBase := "custom_tls/"

	// Create
	var err error
	var cc *CustomCertificate
	record(t, fixtureBase+"create", func(c *Client) {
		cc, err = c.CreateCustomCertificate(&CreateCustomCertificateInput{
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
			c.DeleteCustomCertificate(&DeleteCustomCertificateInput{
				ID: cc.ID,
			})
		})
	}()

	// List
	var lcc []*CustomCertificate
	record(t, fixtureBase+"list", func(c *Client) {
		lcc, err = c.ListCustomCertificates(&ListCustomCertificatesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lcc) < 1 {
		t.Errorf("bad Custom certificates: %v", lcc)
	}

	// Get
	var gcc *CustomCertificate
	record(t, fixtureBase+"get", func(c *Client) {
		gcc, err = c.GetCustomCertificate(&GetCustomCertificateInput{
			ID: cc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if cc.ID != gcc.ID {
		t.Errorf("bad ID: %q (%q)", cc.ID, gcc.ID)
	}

	// Update
	var ucc *CustomCertificate
	record(t, fixtureBase+"update", func(c *Client) {
		ucc, err = c.UpdateCustomCertificate(&UpdateCustomCertificateInput{
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
		err = c.DeleteCustomCertificate(&DeleteCustomCertificateInput{
			ID: cc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateCustomCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/create", func(c *Client) {
		_, err = c.CreateCustomCertificate(&CreateCustomCertificateInput{
			CertBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			Name:     "My certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeleteCustomCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/delete", func(c *Client) {
		err = c.DeleteCustomCertificate(&DeleteCustomCertificateInput{
			ID: "CERTIFICATE_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListCustomCertificates_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/list", func(c *Client) {
		_, err = c.ListCustomCertificates(&ListCustomCertificatesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetCustomCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/get", func(c *Client) {
		_, err = c.GetCustomCertificate(&GetCustomCertificateInput{
			ID: "CERTIFICATE_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_UpdateCustomCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls/update", func(c *Client) {
		_, err = c.UpdateCustomCertificate(&UpdateCustomCertificateInput{
			ID:       "CERTIFICATE_ID",
			CertBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			Name:     "My certificate",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
