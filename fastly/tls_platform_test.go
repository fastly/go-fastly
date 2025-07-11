package fastly

import (
	"context"
	"testing"
)

func TestClient_BulkCertificate(t *testing.T) {
	t.Parallel()

	fixtureBase := "platform_tls/"

	// Create
	var err error
	var bc *BulkCertificate
	Record(t, fixtureBase+"create", func(c *Client) {
		bc, err = c.CreateBulkCertificate(context.TODO(), &CreateBulkCertificateInput{
			CertBlob:          "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			IntermediatesBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			Configurations: []*TLSConfiguration{
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
		Record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteBulkCertificate(context.TODO(), &DeleteBulkCertificateInput{
				ID: bc.ID,
			})
		})
	}()

	// List
	var lbc []*BulkCertificate
	Record(t, fixtureBase+"list", func(c *Client) {
		lbc, err = c.ListBulkCertificates(context.TODO(), &ListBulkCertificatesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lbc) < 1 {
		t.Errorf("bad bulk certificates: %v", lbc)
	}

	// Get
	var gbc *BulkCertificate
	Record(t, fixtureBase+"get", func(c *Client) {
		gbc, err = c.GetBulkCertificate(context.TODO(), &GetBulkCertificateInput{
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
	Record(t, fixtureBase+"update", func(c *Client) {
		ubc, err = c.UpdateBulkCertificate(context.TODO(), &UpdateBulkCertificateInput{
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
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteBulkCertificate(context.TODO(), &DeleteBulkCertificateInput{
			ID: bc.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateBulkCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "platform_tls/create", func(c *Client) {
		_, err = c.CreateBulkCertificate(context.TODO(), &CreateBulkCertificateInput{
			CertBlob:          "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			IntermediatesBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			Configurations: []*TLSConfiguration{
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

func TestClient_DeleteBulkCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "platform_tls/delete", func(c *Client) {
		err = c.DeleteBulkCertificate(context.TODO(), &DeleteBulkCertificateInput{
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
	Record(t, "platform_tls/list", func(c *Client) {
		_, err = c.ListBulkCertificates(context.TODO(), &ListBulkCertificatesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetBulkCertificate_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "platform_tls/get", func(c *Client) {
		_, err = c.GetBulkCertificate(context.TODO(), &GetBulkCertificateInput{
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
	Record(t, "platform_tls/update", func(c *Client) {
		_, err = c.UpdateBulkCertificate(context.TODO(), &UpdateBulkCertificateInput{
			ID:                "CERTIFICATE_ID",
			CertBlob:          "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
			IntermediatesBlob: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
