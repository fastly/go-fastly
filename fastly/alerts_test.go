package fastly

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClient_FastlyAlerts(t *testing.T) {
	t.Parallel()

	testDimensions := map[string][]string{
		"domains": {"example.com", "fastly.com"},
	}
	testEvaluationStrategy := map[string]any{
		"period":    "5m0s",
		"threshold": float64(10),
		"type":      "above_threshold",
	}
	cadi := &CreateAlertDefinitionInput{
		Description:        ToPointer("test description"),
		Dimensions:         testDimensions,
		EvaluationStrategy: testEvaluationStrategy,
		IntegrationIDs:     []string{},
		Metric:             ToPointer("status_5xx"),
		Name:               ToPointer("test name"),
		ServiceID:          ToPointer(testServiceID),
		Source:             ToPointer("domains"),
	}

	// Test
	var err error
	record(t, "alerts/test_alert_definition", func(c *Client) {
		err = c.TestAlertDefinition(&TestAlertDefinitionInput{
			CreateAlertDefinitionInput: *cadi,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create
	var ad *AlertDefinition
	record(t, "alerts/create_alert_definition", func(c *Client) {
		ad, err = c.CreateAlertDefinition(cadi)

	})
	if err != nil {
		t.Fatal(err)
	}
	// Ensure deleted
	defer func() {
		record(t, "alerts/cleanup_alert_definition", func(c *Client) {
			err = c.DeleteAlertDefinition(&DeleteAlertDefinitionInput{
				ID: &ad.ID,
			})
		})
	}()

	if ad.Description != "test description" {
		t.Errorf("bad description: %v", ad.Description)
	}

	if ad.Metric != "status_5xx" {
		t.Errorf("bad metric: %v", ad.Metric)
	}

	if ad.Name != "test name" {
		t.Errorf("bad name: %v", ad.Name)
	}

	if ad.ServiceID != testServiceID {
		t.Errorf("bad service_id: %v", ad.ServiceID)
	}

	if ad.Source != "domains" {
		t.Errorf("bad source: %v", ad.Source)
	}

	if diff := cmp.Diff(testDimensions, ad.Dimensions); diff != "" {
		t.Errorf("bad dimensions: diff -want +got\n%v", diff)
	}

	if diff := cmp.Diff(testEvaluationStrategy, ad.EvaluationStrategy); diff != "" {
		t.Errorf("bad evaluation_strategy: diff -want +got\n%v", diff)
	}

	// List Definitions
	var adr *AlertDefinitionsResponse
	record(t, "alerts/list_alert_definitions", func(c *Client) {
		adr, err = c.ListAlertDefinitions(&ListAlertDefinitionsInput{
			Cursor:    ToPointer(""),
			Limit:     ToPointer(10),
			Name:      ToPointer(ad.Name),
			ServiceID: ToPointer(testServiceID),
			Sort:      ToPointer("name"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(adr.Data) < 1 {
		t.Errorf("bad alert definitions: %v", adr)
	}

	// Get
	var gad *AlertDefinition
	record(t, "alerts/get_alert_definition", func(c *Client) {
		gad, err = c.GetAlertDefinition(&GetAlertDefinitionInput{
			ID: &ad.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ad.Name != gad.Name {
		t.Errorf("bad name: %q (%q)", ad.Name, gad.Name)
	}

	// Update
	var uad *AlertDefinition
	record(t, "alerts/update_alert_definition", func(c *Client) {
		uad, err = c.UpdateAlertDefinition(&UpdateAlertDefinitionInput{
			Description:        ToPointer("test description"),
			Dimensions:         testDimensions,
			EvaluationStrategy: testEvaluationStrategy,
			ID:                 ToPointer(ad.ID),
			IntegrationIDs:     []string{},
			Metric:             ToPointer("status_5xx"),
			Name:               ToPointer("test name updated"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uad.Name != "test name updated" {
		t.Errorf("bad name: %v", uad.Name)
	}

	// Delete
	record(t, "alerts/delete_alert_definition", func(c *Client) {
		err = c.DeleteAlertDefinition(&DeleteAlertDefinitionInput{
			ID: &ad.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// List History
	record(t, "alerts/list_alert_history", func(c *Client) {
		_, err = c.ListAlertHistory(&ListAlertHistoryInput{
			After:        ToPointer("2006-01-02T15:04:05Z"),
			Before:       ToPointer("2056-01-02T15:04:05Z"),
			Cursor:       ToPointer(""),
			DefinitionID: ToPointer(ad.ID),
			Limit:        ToPointer(10),
			ServiceID:    ToPointer(testServiceID),
			Sort:         ToPointer("-start"),
			Status:       ToPointer(""),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetAlertDefinition_validation(t *testing.T) {
	var err error
	_, err = testClient.GetAlertDefinition(&GetAlertDefinitionInput{
		ID: nil,
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateAlertDefinition_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateAlertDefinition(&UpdateAlertDefinitionInput{
		ID: nil,
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteAlertDefinition_validation(t *testing.T) {
	err := testClient.DeleteAlertDefinition(&DeleteAlertDefinitionInput{
		ID: nil,
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
