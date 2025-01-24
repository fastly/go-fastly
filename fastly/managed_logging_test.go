package fastly

import (
	"errors"
	"testing"
)

func TestClient_ManagedLogging(t *testing.T) {
	t.Parallel()

	var err error

	// Create
	Record(t, "managed_logging/create", func(c *Client) {
		_, err = c.CreateManagedLogging(&CreateManagedLoggingInput{
			ServiceID: TestDeliveryServiceID,
			Kind:      ManagedLoggingInstanceOutput,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that enabling managed logging on a service with it already
	// enabled results in a 409.
	Record(t, "managed_logging/recreate", func(c *Client) {
		_, err = c.CreateManagedLogging(&CreateManagedLoggingInput{
			ServiceID: TestDeliveryServiceID,
			Kind:      ManagedLoggingInstanceOutput,
		})
	})
	if !errors.Is(err, ErrManagedLoggingEnabled) {
		t.Errorf("unexpected error: %s", err)
	}

	// Delete
	Record(t, "managed_logging/delete", func(c *Client) {
		err = c.DeleteManagedLogging(&DeleteManagedLoggingInput{
			ServiceID: TestDeliveryServiceID,
			Kind:      ManagedLoggingInstanceOutput,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateManagedLogging_validation(t *testing.T) {
	_, err := TestClient.CreateManagedLogging(&CreateManagedLoggingInput{
		ServiceID: "",
		Kind:      ManagedLoggingInstanceOutput,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = TestClient.CreateManagedLogging(&CreateManagedLoggingInput{
		ServiceID: TestDeliveryServiceID,
	})
	if !errors.Is(err, ErrMissingKind) {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = TestClient.CreateManagedLogging(&CreateManagedLoggingInput{
		ServiceID: TestDeliveryServiceID,
		Kind:      999,
	})
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestClient_DeleteManagedLogging_validation(t *testing.T) {
	err := TestClient.DeleteManagedLogging(&DeleteManagedLoggingInput{
		ServiceID: "",
		Kind:      ManagedLoggingInstanceOutput,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("unexpected error: %s", err)
	}

	err = TestClient.DeleteManagedLogging(&DeleteManagedLoggingInput{
		ServiceID: TestDeliveryServiceID,
	})
	if !errors.Is(err, ErrMissingKind) {
		t.Errorf("unexpected error: %s", err)
	}

	err = TestClient.DeleteManagedLogging(&DeleteManagedLoggingInput{
		ServiceID: TestDeliveryServiceID,
		Kind:      999,
	})
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("unexpected error: %s", err)
	}
}
