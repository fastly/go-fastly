package fastly

import "testing"

func TestClient_PrivateKey(t *testing.T) {
	t.Parallel()

	fixtureBase := "platform_tls/"

	// Create
	var err error
	var pk *PrivateKey
	record(t, fixtureBase+"create", func(c *Client) {
		pk, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
			Name: "My private key",
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

func TestClient_BulkCertificate(t *testing.T) {
	t.Parallel()

	fixtureBase := "platform_tls/"

	// Create
	var err error
	var bc *BulkCertificate
	record(t, fixtureBase+"create_bulk", func(c *Client) {
		bc, err = c.CreateBulkCertificate(&CreateBulkCertificateInput{
			CertBlob:          "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			IntermediatesBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			TLSConfigurations: []*TLSConfiguration{
				{
					ID: "TLS_CONFIGURATION_ID",
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup_bulk", func(c *Client) {
			c.DeleteBulkCertificate(&DeleteBulkCertificateInput{
				ID: bc.ID,
			})
		})
	}()

	// List
	var lbc []*BulkCertificate
	record(t, fixtureBase+"list_bulk", func(c *Client) {
		lbc, err = c.ListBulkCertificates(&ListBulkCertificatesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lbc) < 1 {
		t.Errorf("bad bulk certificates: %v", lbc)
	}

	// Get
	var gbc *BulkCertificate
	record(t, fixtureBase+"get_bulk", func(c *Client) {
		gbc, err = c.GetBulkCertificate(&GetBulkCertificateInput{
			ID: bc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if bc.ID != gbc.ID {
		t.Errorf("bad ID: %q (%q)", bc.ID, gbc.ID)
	}

	// Update
	var ubc *BulkCertificate
	record(t, fixtureBase+"update_bulk", func(c *Client) {
		ubc, err = c.UpdateBulkCertificate(&UpdateBulkCertificateInput{
			ID:                "CERTIFICATE_ID",
			CertBlob:          "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			IntermediatesBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if bc.ID != ubc.ID {
		t.Errorf("bad ID: %q (%q)", bc.ID, ubc.ID)
	}

	// Delete
	record(t, fixtureBase+"delete_bulk", func(c *Client) {
		err = c.DeleteBulkCertificate(&DeleteBulkCertificateInput{
			ID: bc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPrivateKeys_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "platform_tls/list", func(c *Client) {
		_, err = c.ListPrivateKeys(&ListPrivateKeysInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetPrivateKey_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "platform_tls/get", func(c *Client) {
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
	record(t, "platform_tls/create", func(c *Client) {
		_, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
			Name: "My private key",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateBulkCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "platform_tls/create_bulk", func(c *Client) {
		_, err = c.CreateBulkCertificate(&CreateBulkCertificateInput{
			CertBlob:          "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			IntermediatesBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			TLSConfigurations: []*TLSConfiguration{
				{
					ID: "TLS_CONFIGURATION_ID",
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeletePrivateKey_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "platform_tls/delete", func(c *Client) {
		err = c.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: "PRIVATE_KEY_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeleteBulkCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "platform_tls/delete_bulk", func(c *Client) {
		err = c.DeleteBulkCertificate(&DeleteBulkCertificateInput{
			ID: "CERTIFICATE_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListBulkCertificates_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "platform_tls/list_bulk", func(c *Client) {
		_, err = c.ListBulkCertificates(&ListBulkCertificatesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetBulkCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "platform_tls/get_bulk", func(c *Client) {
		_, err = c.GetBulkCertificate(&GetBulkCertificateInput{
			ID: "CERTIFICATE_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_UpdateBulkCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "platform_tls/update_bulk", func(c *Client) {
		_, err = c.UpdateBulkCertificate(&UpdateBulkCertificateInput{
			ID:                "CERTIFICATE_ID",
			CertBlob:          "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			IntermediatesBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
