package fastly

import (
	"testing"
)

func TestClient_CustomTLSConfiguration(t *testing.T) {
	t.Parallel()

	fixtureBase := "custom_tls_configuration/"

	var err error
	conID := "TLS_CONFIGURATION_ID"

	// Get
	var gcon *CustomTLSConfiguration
	record(t, fixtureBase+"get", func(c *Client) {
		gcon, err = c.GetCustomTLSConfiguration(&GetCustomTLSConfigurationInput{
			ID: conID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if conID != gcon.ID {
		t.Errorf("bad ID: %q (%q)", conID, gcon.ID)
	}

	// List
	var lcon []*CustomTLSConfiguration
	record(t, fixtureBase+"list", func(c *Client) {
		lcon, err = c.ListCustomTLSConfigurations(&ListCustomTLSConfigurationsInput{})
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
	record(t, fixtureBase+"update", func(c *Client) {
		ucon, err = c.UpdateCustomTLSConfiguration(&UpdateCustomTLSConfigurationInput{
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

}

func TestClient_ListCustomTLSConfigurations_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_configuration/list", func(c *Client) {
		_, err = c.ListCustomTLSConfigurations(&ListCustomTLSConfigurationsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetCustomTLSConfiguration_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_configuration/get", func(c *Client) {
		_, err = c.GetCustomTLSConfiguration(&GetCustomTLSConfigurationInput{
			ID: "TLS_CONFIGURATION_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.GetCustomTLSConfiguration(&GetCustomTLSConfigurationInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCustomTLSConfiguration_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "custom_tls_configuration/update", func(c *Client) {
		_, err = c.UpdateCustomTLSConfiguration(&UpdateCustomTLSConfigurationInput{
			ID:   "TLS_CONFIGURATION_ID",
			Name: "My configuration v2",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.UpdateCustomTLSConfiguration(&UpdateCustomTLSConfigurationInput{
		Name: "My configuration v2",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCustomTLSConfiguration(&UpdateCustomTLSConfigurationInput{
		ID: "CONFIGURATION_ID",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
