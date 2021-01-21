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
			Certificate:   &CustomTLSCertificate{ID: "CERTIFICATE_ID"},
			Configuration: &TLSConfiguration{ID: "CONFIGURATION_ID"},
			Domain:        &TLSDomain{ID: "DOMAIN_NAME"},
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
	if lta[0].Certificate == nil {
		t.Errorf("TLS certificate relation should not be nil: %v", lta)
	}
	if lta[0].Certificate.ID != ta.Certificate.ID {
		t.Errorf("bad Certificate ID: %q (%q)", lta[0].Certificate.ID, ta.Certificate.ID)
	}
	if lta[0].Configuration == nil {
		t.Errorf("TLS Configuration relation should not be nil: %v", lta)
	}
	if lta[0].Configuration.ID != ta.Configuration.ID {
		t.Errorf("bad Configuration ID: %q (%q)", lta[0].Configuration.ID, ta.Configuration.ID)
	}
	if lta[0].Domain == nil {
		t.Errorf("TLS domain relation should not be nil: %v", lta)
	}
	if lta[0].Domain.ID != ta.Domain.ID {
		t.Errorf("bad Domain ID: %q (%q)", lta[0].Domain.ID, ta.Domain.ID)
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
			ID:          "ACTIVATION_ID",
			Certificate: &CustomTLSCertificate{},
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
			Certificate:   &CustomTLSCertificate{ID: "CERTIFICATE_ID"},
			Configuration: &TLSConfiguration{ID: "CONFIGURATION_ID"},
			Domain:        &TLSDomain{ID: "DOMAIN_NAME"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.CreateTLSActivation(&CreateTLSActivationInput{
		Configuration: &TLSConfiguration{ID: "CONFIGURATION_ID"},
		Domain:        &TLSDomain{ID: "DOMAIN_NAME"},
	})
	if err != ErrMissingTLSCertificate {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateTLSActivation(&CreateTLSActivationInput{
		Certificate:   &CustomTLSCertificate{ID: "CERTIFICATE_ID"},
		Configuration: &TLSConfiguration{ID: "CONFIGURATION_ID"},
	})
	if err != ErrMissingTLSDomain {
		t.Errorf("bad error: %s", err)
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

	err = testClient.DeleteTLSActivation(&DeleteTLSActivationInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
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

	_, err = testClient.GetTLSActivation(&GetTLSActivationInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateTLSActivation_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_activation/update", func(c *Client) {
		_, err = c.UpdateTLSActivation(&UpdateTLSActivationInput{
			ID:          "ACTIVATION_ID",
			Certificate: &CustomTLSCertificate{ID: "CERTIFICATE_ID"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.UpdateTLSActivation(&UpdateTLSActivationInput{
		ID: "ACTIVATION_ID",
	})
	if err != ErrMissingTLSCertificate {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateTLSActivation(&UpdateTLSActivationInput{
		Certificate: &CustomTLSCertificate{ID: "CERTIFICATE_ID"},
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
