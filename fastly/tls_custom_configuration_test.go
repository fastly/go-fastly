package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_CustomTLSConfiguration(t *testing.T) {
	t.Parallel()

	fixtureBase := "custom_tls_configuration/"

	var err error
	conID := "TLS_CONFIGURATION_ID"

	// Create
	var ccon *CustomTLSConfiguration
	Record(t, fixtureBase+"create", func(c *Client) {
		ccon, err = c.CreateCustomTLSConfiguration(context.TODO(), &CreateCustomTLSConfigurationInput{
			Name:          "My configuration",
			HTTPProtocols: []string{"http/1.1", "http/2"},
			TLSProtocols:  []string{"1.2", "1.3"},
			Vipspace:      ToPointer("myvipspace"),
			DefaultCertificate: &TLSCertificateRef{
				ID: "1TTgJM5iiNcHi8V5cwCYk0",
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if conID != ccon.ID {
		t.Errorf("bad ID: %q (%q)", conID, ccon.ID)
	}
	if len(ccon.TLS12CipherSuiteProfile) == 0 {
		t.Errorf("expected TLS12CipherSuiteProfile to be populated")
	}
	if len(ccon.TLS13CipherSuiteProfile) == 0 {
		t.Errorf("expected TLS13CipherSuiteProfile to be populated")
	}
	if ccon.DefaultCertificate == nil || ccon.DefaultCertificate.ID != "1TTgJM5iiNcHi8V5cwCYk0" {
		t.Errorf("bad DefaultCertificate: %+v", ccon.DefaultCertificate)
	}
	if ccon.Vipspace == nil || *ccon.Vipspace != "myvipspace" {
		t.Errorf("bad Vipspace: %+v", ccon.Vipspace)
	}

	// Get
	var gcon *CustomTLSConfiguration
	Record(t, fixtureBase+"get", func(c *Client) {
		gcon, err = c.GetCustomTLSConfiguration(context.TODO(), &GetCustomTLSConfigurationInput{
			ID: conID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if conID != gcon.ID {
		t.Errorf("bad ID: %q (%q)", conID, gcon.ID)
	}
	if gcon.Vipspace == nil || *gcon.Vipspace != "myvipspace" {
		t.Errorf("bad Vipspace: %v", gcon.Vipspace)
	}
	if len(gcon.TLS12CipherSuiteProfile) == 0 {
		t.Errorf("expected TLS12CipherSuiteProfile to be populated")
	}
	if len(gcon.TLS13CipherSuiteProfile) == 0 {
		t.Errorf("expected TLS13CipherSuiteProfile to be populated")
	}
	if gcon.DefaultCertificate == nil || gcon.DefaultCertificate.ID != "1TTgJM5iiNcHi8V5cwCYk0" {
		t.Errorf("bad DefaultCertificate: %+v", gcon.DefaultCertificate)
	}

	// List
	var lcon []*CustomTLSConfiguration
	Record(t, fixtureBase+"list", func(c *Client) {
		lcon, err = c.ListCustomTLSConfigurations(context.TODO(), &ListCustomTLSConfigurationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lcon) < 1 {
		t.Errorf("bad tls configurations: %v", lcon)
	}

	// Update
	var ucon *CustomTLSConfiguration
	newName := "My configuration v2"
	Record(t, fixtureBase+"update", func(c *Client) {
		ucon, err = c.UpdateCustomTLSConfiguration(context.TODO(), &UpdateCustomTLSConfigurationInput{
			ID:   "TLS_CONFIGURATION_ID",
			Name: newName,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if conID != ucon.ID {
		t.Errorf("bad ID: %q (%q)", conID, ucon.ID)
	}
	if ucon.Name != newName {
		t.Errorf("bad Name: %q (%q)", newName, ucon.Name)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteCustomTLSConfiguration(context.TODO(), &DeleteCustomTLSConfigurationInput{
			ID: conID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListCustomTLSConfigurations_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "custom_tls_configuration/list", func(c *Client) {
		_, err = c.ListCustomTLSConfigurations(context.TODO(), &ListCustomTLSConfigurationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetCustomTLSConfiguration_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "custom_tls_configuration/get", func(c *Client) {
		_, err = c.GetCustomTLSConfiguration(context.TODO(), &GetCustomTLSConfigurationInput{
			ID: "TLS_CONFIGURATION_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = TestClient.GetCustomTLSConfiguration(context.TODO(), &GetCustomTLSConfigurationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCustomTLSConfiguration_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, "custom_tls_configuration/update", func(c *Client) {
		_, err = c.UpdateCustomTLSConfiguration(context.TODO(), &UpdateCustomTLSConfigurationInput{
			ID:   "TLS_CONFIGURATION_ID",
			Name: "My configuration v2",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = TestClient.UpdateCustomTLSConfiguration(context.TODO(), &UpdateCustomTLSConfigurationInput{
		Name: "My configuration v2",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateCustomTLSConfiguration(context.TODO(), &UpdateCustomTLSConfigurationInput{
		ID: "CONFIGURATION_ID",
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateCustomTLSConfiguration_validation(t *testing.T) {
	t.Parallel()

	_, err := TestClient.CreateCustomTLSConfiguration(context.TODO(), &CreateCustomTLSConfigurationInput{})
	if !errors.Is(err, ErrMissingHTTPProtocols) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCustomTLSConfiguration_validation(t *testing.T) {
	t.Parallel()

	err := TestClient.DeleteCustomTLSConfiguration(context.TODO(), &DeleteCustomTLSConfigurationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}
