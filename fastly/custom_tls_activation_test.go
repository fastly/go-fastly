package fastly

import (
	"testing"
)

func TestClient_TLSActivation(t *testing.T) {
	t.Parallel()

	fixtureBase := "custom_tls_activation/"

	// Create
	var err error
	var ta *TLSActivation
	record(t, fixtureBase+"create", func(c *Client) {
		ta, err = c.CreateTLSActivation(&CreateTLSActivationInput{
			TLSCertificate:   &TLSCertificate{ID: "CERTIFICATE_ID"},
			TLSConfiguration: &TLSConfiguration{ID: "CONFIGURATION_ID"},
			TLSDomain:        &TLSDomain{ID: "DOMAIN_NAME"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteTLSActivation(&DeleteTLSActivationInput{
				ID: ta.ID,
			})
		})
	}()

	// List
	var lta []*TLSActivation
	record(t, fixtureBase+"list", func(c *Client) {
		lta, err = c.ListTLSActivations(&ListTLSActivationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lta) < 1 {
		t.Errorf("bad TLS activations: %v", lta)
	}
	if lta[0].TLSCertificate == nil {
		t.Errorf("TLS certificate relation should not be nil: %v", lta)
	}
	if lta[0].TLSConfiguration == nil {
		t.Errorf("TLS configuration relation should not be nil: %v", lta)
	}
	if lta[0].TLSDomain == nil {
		t.Errorf("TLS domain relation should not be nil: %v", lta)
	}

	// Get
	var gta *TLSActivation
	record(t, fixtureBase+"get", func(c *Client) {
		gta, err = c.GetTLSActivation(&GetTLSActivationInput{
			ID: ta.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ta.ID != gta.ID {
		t.Errorf("bad ID: %q (%q)", ta.ID, gta.ID)
	}

	// Update
	var uta *TLSActivation
	record(t, fixtureBase+"update", func(c *Client) {
		uta, err = c.UpdateTLSActivation(&UpdateTLSActivationInput{
			ID:             "ACTIVATION_ID",
			TLSCertificate: &TLSCertificate{},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ta.ID != uta.ID {
		t.Errorf("bad ID: %q (%q)", ta.ID, uta.ID)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteTLSActivation(&DeleteTLSActivationInput{
			ID: ta.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateTLSActivation_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_activation/create", func(c *Client) {
		_, err = c.CreateTLSActivation(&CreateTLSActivationInput{
			TLSCertificate:   &TLSCertificate{},
			TLSConfiguration: &TLSConfiguration{},
			TLSDomain:        &TLSDomain{},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeleteTLSActivation_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_activation/delete", func(c *Client) {
		err = c.DeleteTLSActivation(&DeleteTLSActivationInput{
			ID: "ACTIVATION_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListTLSActivations_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_activation/list", func(c *Client) {
		_, err = c.ListTLSActivations(&ListTLSActivationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetTLSActivation_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_activation/get", func(c *Client) {
		_, err = c.GetTLSActivation(&GetTLSActivationInput{
			ID: "ACTIVATION_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_UpdateTLSActivation_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_activation/update", func(c *Client) {
		_, err = c.UpdateTLSActivation(&UpdateTLSActivationInput{
			ID:             "ACTIVATION_ID",
			TLSCertificate: &TLSCertificate{ID: "CERTIFICATE_ID"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
