package fastly

import (
	"errors"
	"testing"
)

func TestClient_Notifications(t *testing.T) {
	t.Parallel()
	var err error

	// Get integration types
	var its *[]IntegrationType
	Record(t, "notifications/get_integration_types", func(c *Client) {
		its, err = c.GetIntegrationTypes()
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
		cir, err = c.CreateIntegration(cii)
	})
	// Ensure integration deleted
	defer func() {
		Record(t, "notifications/cleanup_integration", func(c *Client) {
			err = c.DeleteIntegration(&DeleteIntegrationInput{
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
		sir, err = c.SearchIntegrations(&SearchIntegrationsInput{
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
		gi, err = c.GetIntegration(&GetIntegrationInput{
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
		err = c.CreateMailinglistConfirmation(&CreateMailinglistConfirmationInput{
			Email: ToPointer("noreply@fastly.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Update integration
	Record(t, "notifications/update_integration", func(c *Client) {
		err = c.UpdateIntegration(&UpdateIntegrationInput{
			Config: map[string]string{
				"url": "https://foo.com/bar",
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
		rwskr, err = c.RotateWebhookSigningKey(&RotateWebhookSigningKeyInput{
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
		gwskr, err = c.GetWebhookSigningKey(&GetWebhookSigningKeyInput{
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
		err = c.DeleteIntegration(&DeleteIntegrationInput{
			ID: *gi.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetIntegration_validation(t *testing.T) {
	_, err := TestClient.GetIntegration(&GetIntegrationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateIntegration_validation(t *testing.T) {
	err := TestClient.UpdateIntegration(&UpdateIntegrationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteIntegration_validation(t *testing.T) {
	err := TestClient.DeleteIntegration(&DeleteIntegrationInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetWebhookSigningKey_validation(t *testing.T) {
	_, err := TestClient.GetWebhookSigningKey(&GetWebhookSigningKeyInput{})
	if !errors.Is(err, ErrMissingIntegrationID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_RotateWebhookSigningKey_validation(t *testing.T) {
	_, err := TestClient.RotateWebhookSigningKey(&RotateWebhookSigningKeyInput{})
	if !errors.Is(err, ErrMissingIntegrationID) {
		t.Errorf("bad error: %s", err)
	}
}
