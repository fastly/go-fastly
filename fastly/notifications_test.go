package fastly

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestClient_Notifications(t *testing.T) {
	t.Parallel()
	var err error

	// Get integration types
	var its *[]IntegrationType
	Record(t, "notifications/get_integration_types", func(c *Client) {
		its, err = c.GetIntegrationTypes(context.TODO())
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(*its) < 6 {
		t.Errorf("missing integration types, %v", its)
	}

	cii := &CreateIntegrationInput{
		Config: map[string]string{
			"address": "noreply@fastly.com",
		},
		Description: ToPointer("test description"),
		Name:        ToPointer("test name"),
		Type:        ToPointer("mailinglist"),
	}

	// Create integration
	var cir *CreateIntegrationResponse
	Record(t, "notifications/create_integration", func(c *Client) {
		cir, err = c.CreateIntegration(context.TODO(), cii)
	})
	// Ensure integration deleted
	defer func() {
		Record(t, "notifications/cleanup_integration", func(c *Client) {
			err = c.DeleteIntegration(context.TODO(), &DeleteIntegrationInput{
				ID: *cir.ID,
			})
		})
	}()
	if cir.ID == nil {
		t.Errorf("missing id")
	}

	// Search integrations
	var sir *SearchIntegrationsResponse
	Record(t, "notifications/search_integrations", func(c *Client) {
		sir, err = c.SearchIntegrations(context.TODO(), &SearchIntegrationsInput{
			Cursor: ToPointer(""),
			Limit:  ToPointer(3),
			Sort:   ToPointer("-created_at"),
			Type:   ToPointer("mailinglist"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sir.Data) < 1 {
		t.Errorf("bad data: %v", sir.Data)
	}
	if *sir.Meta.Limit != 3 {
		t.Errorf("bad meta limit: %v", sir.Meta.Limit)
	}
	if *sir.Meta.Sort != "-created_at" {
		t.Errorf("bad meta sort, %v", sir.Meta.Sort)
	}
	if *sir.Meta.Total < 1 {
		t.Errorf("bad meta total, %v", sir.Meta.Total)
	}
	if *sir.Meta.Type != "mailinglist" {
		t.Errorf("bad meta type, %v", sir.Meta.Type)
	}

	// Get integration
	var gi *Integration
	Record(t, "notifications/get_integration", func(c *Client) {
		gi, err = c.GetIntegration(context.TODO(), &GetIntegrationInput{
			ID: *cir.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gi.CreatedAt.IsZero() {
		t.Errorf("missing created at")
	}
	if *gi.Description != *cii.Description {
		t.Errorf("bad description: %q (%q)", *gi.Description, *cii.Description)
	}
	if *gi.ID != *cir.ID {
		t.Errorf("bad id: %q (%q)", *gi.ID, *cir.ID)
	}
	if *gi.Name != *cii.Name {
		t.Errorf("bad name: %q (%q)", *gi.Name, *cii.Name)
	}
	if *gi.Type != *cii.Type {
		t.Errorf("bad type: %q (%q)", *gi.Type, *cii.Type)
	}
	if gi.UpdatedAt.IsZero() {
		t.Errorf("missing updated at")
	}

	// Create mailinglist integration confirmation
	Record(t, "notifications/create_mailinglist_confirmation", func(c *Client) {
		err = c.CreateMailinglistConfirmation(context.TODO(), &CreateMailinglistConfirmationInput{
			Email: ToPointer("noreply@fastly.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Update integration
	Record(t, "notifications/update_integration", func(c *Client) {
		err = c.UpdateIntegration(context.TODO(), &UpdateIntegrationInput{
			Config: map[string]string{
				"webhook": "https://foo.com/bar",
			},
			Description: ToPointer("test description updated"),
			ID:          *gi.ID,
			Name:        ToPointer("test name updated"),
			Type:        ToPointer("webhook"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Rotate webhook integration signing key
	var rwskr *WebhookSigningKeyResponse
	Record(t, "notifications/rotate_webhook_signing_key", func(c *Client) {
		rwskr, err = c.RotateWebhookSigningKey(context.TODO(), &RotateWebhookSigningKeyInput{
			IntegrationID: *gi.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if rwskr.SigningKey == nil {
		t.Errorf("rotate missing signing key")
	}

	// Get webhook integration signing key
	var gwskr *WebhookSigningKeyResponse
	Record(t, "notifications/get_webhook_signing_key", func(c *Client) {
		gwskr, err = c.GetWebhookSigningKey(context.TODO(), &GetWebhookSigningKeyInput{
			IntegrationID: *gi.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if gwskr.SigningKey == nil {
		t.Errorf("get missing signing key")
	}

	// Delete integration
	Record(t, "notifications/delete_integration", func(c *Client) {
		err = c.DeleteIntegration(context.TODO(), &DeleteIntegrationInput{
			ID: *gi.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create, get, and delete an integration for each of the newer
	// integration types, using their typed config helpers.
	newTypeIntegrations := []struct {
		fixture string
		typ     string
		config  map[string]string
	}{
		{
			fixture: "datadog",
			typ:     IntegrationTypeDatadog,
			config:  DatadogConfig{APIKey: "abc123def456", Site: "datadoghq.eu"}.ToMap(),
		},
		{
			fixture: "jiraissue",
			typ:     IntegrationTypeJiraIssue,
			config: JiraIssueConfig{ //nolint: gosec
				BaseURL:    "https://your-org.atlassian.net",
				Username:   "user@example.com",
				Token:      "AbCdEfGhIjKlMnOpQrStUvWx",
				ProjectKey: "PROJ",
				IssueType:  "Bug",
			}.ToMap(),
		},
		{
			fixture: "jsm",
			typ:     IntegrationTypeJSM,
			config:  JSMConfig{APIKey: "g4ff854d-a14c-46a8-b8f0-0960774319dd"}.ToMap(), //nolint: gosec
		},
		{
			fixture: "opsgenie",
			typ:     IntegrationTypeOpsGenie,
			config:  OpsGenieConfig{APIKey: "a1bb854d-b24c-46a8-c9f0-1960774319ee"}.ToMap(), //nolint: gosec
		},
		{
			fixture: "splunkoncall",
			typ:     IntegrationTypeSplunkOnCall,
			config:  SplunkOnCallConfig{URL: "https://alert.victorops.com/integrations/generic/20131114/alert/XXXX"}.ToMap(),
		},
	}

	for _, nti := range newTypeIntegrations {
		// Create integration
		var newCir *CreateIntegrationResponse
		Record(t, "notifications/create_"+nti.fixture+"_integration", func(c *Client) {
			newCir, err = c.CreateIntegration(context.TODO(), &CreateIntegrationInput{
				Config:      nti.config,
				Description: ToPointer("test " + nti.fixture + " description"),
				Name:        ToPointer("test " + nti.fixture + " name"),
				Type:        ToPointer(nti.typ),
			})
		})
		if err != nil {
			t.Fatal(err)
		}
		if newCir == nil || newCir.ID == nil {
			t.Fatalf("missing response or integration id for %s integration", nti.fixture)
		}

		// Get integration
		var newGi *Integration
		Record(t, "notifications/get_"+nti.fixture+"_integration", func(c *Client) {
			newGi, err = c.GetIntegration(context.TODO(), &GetIntegrationInput{
				ID: *newCir.ID,
			})
		})
		if err != nil {
			t.Fatal(err)
		}
		if *newGi.Type != nti.typ {
			t.Errorf("bad type for %s integration: %q", nti.fixture, *newGi.Type)
		}
		if *newGi.ID != *newCir.ID {
			t.Errorf("bad id for %s integration: %q (%q)", nti.fixture, *newGi.ID, *newCir.ID)
		}

		// Delete integration
		Record(t, "notifications/delete_"+nti.fixture+"_integration", func(c *Client) {
			err = c.DeleteIntegration(context.TODO(), &DeleteIntegrationInput{
				ID: *newCir.ID,
			})
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestClient_GetIntegration_validation(t *testing.T) {
	_, err := TestClient.GetIntegration(context.TODO(), &GetIntegrationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateIntegration_validation(t *testing.T) {
	err := TestClient.UpdateIntegration(context.TODO(), &UpdateIntegrationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteIntegration_validation(t *testing.T) {
	err := TestClient.DeleteIntegration(context.TODO(), &DeleteIntegrationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestDatadogConfig_ToMap(t *testing.T) {
	got := DatadogConfig{APIKey: "abc123", Site: "datadoghq.eu"}.ToMap()
	want := map[string]string{"apikey": "abc123", "site": "datadoghq.eu"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("bad map: got %v, want %v", got, want)
	}

	// Site is optional and should be omitted when empty.
	got = DatadogConfig{APIKey: "abc123"}.ToMap()
	want = map[string]string{"apikey": "abc123"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("bad map: got %v, want %v", got, want)
	}
}

func TestJiraIssueConfig_ToMap(t *testing.T) {
	got := JiraIssueConfig{ //nolint: gosec
		BaseURL:    "https://your-org.atlassian.net",
		Username:   "user@example.com",
		Token:      "AbCdEfGhIjKlMnOpQrStUvWx",
		ProjectKey: "PROJ",
		IssueType:  "Bug",
	}.ToMap()
	want := map[string]string{ //nolint: gosec
		"baseurl":    "https://your-org.atlassian.net",
		"username":   "user@example.com",
		"token":      "AbCdEfGhIjKlMnOpQrStUvWx",
		"projectkey": "PROJ",
		"issuetype":  "Bug",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("bad map: got %v, want %v", got, want)
	}
}

func TestJSMConfig_ToMap(t *testing.T) {
	got := JSMConfig{APIKey: "g4ff854d-a14c-46a8-b8f0-0960774319dd"}.ToMap()    //nolint: gosec
	want := map[string]string{"apikey": "g4ff854d-a14c-46a8-b8f0-0960774319dd"} //nolint: gosec
	if !reflect.DeepEqual(got, want) {
		t.Errorf("bad map: got %v, want %v", got, want)
	}
}

func TestOpsGenieConfig_ToMap(t *testing.T) {
	got := OpsGenieConfig{APIKey: "a1bb854d-b24c-46a8-c9f0-1960774319ee"}.ToMap() //nolint: gosec
	want := map[string]string{"apikey": "a1bb854d-b24c-46a8-c9f0-1960774319ee"}   //nolint: gosec
	if !reflect.DeepEqual(got, want) {
		t.Errorf("bad map: got %v, want %v", got, want)
	}
}

func TestSplunkOnCallConfig_ToMap(t *testing.T) {
	got := SplunkOnCallConfig{URL: "https://alert.victorops.com/integrations/generic/20131114/alert/XXXX"}.ToMap()
	want := map[string]string{"url": "https://alert.victorops.com/integrations/generic/20131114/alert/XXXX"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("bad map: got %v, want %v", got, want)
	}
}

func TestClient_GetWebhookSigningKey_validation(t *testing.T) {
	_, err := TestClient.GetWebhookSigningKey(context.TODO(), &GetWebhookSigningKeyInput{})
	if !errors.Is(err, ErrMissingIntegrationID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_RotateWebhookSigningKey_validation(t *testing.T) {
	_, err := TestClient.RotateWebhookSigningKey(context.TODO(), &RotateWebhookSigningKeyInput{})
	if !errors.Is(err, ErrMissingIntegrationID) {
		t.Errorf("bad error: %s", err)
	}
}
