package fastly

import "testing"

func TestClient_ManagedLogging(t *testing.T) {
	t.Parallel()

	var err error

	// Create
	record(t, "managed_logging/create", func(c *Client) {
		_, err = c.CreateManagedLogging(&CreateManagedLoggingInput{
			ServiceID: testServiceID,
			Kind:      ManagedLoggingInstanceOutput,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that enabling managed logging on a service with it already
	// enabled results in a 409.
	record(t, "managed_logging/recreate", func(c *Client) {
		_, err = c.CreateManagedLogging(&CreateManagedLoggingInput{
			ServiceID: testServiceID,
			Kind:      ManagedLoggingInstanceOutput,
		})
	})
	if err != ErrManagedLoggingEnabled {
		t.Errorf("unexpected error: %s", err)
	}

	// Delete
	record(t, "managed_logging/delete", func(c *Client) {
		err = c.DeleteManagedLogging(&DeleteManagedLoggingInput{
			ServiceID: testServiceID,
			Kind:      ManagedLoggingInstanceOutput,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateManagedLogging_validation(t *testing.T) {
	_, err := testClient.CreateManagedLogging(&CreateManagedLoggingInput{
		ServiceID: "",
		Kind:      ManagedLoggingInstanceOutput,
	})
	if err != ErrMissingServiceID {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = testClient.CreateManagedLogging(&CreateManagedLoggingInput{
		ServiceID: testServiceID,
	})
	if err != ErrMissingKind {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = testClient.CreateManagedLogging(&CreateManagedLoggingInput{
		ServiceID: testServiceID,
		Kind:      999,
	})
	if err != ErrNotImplemented {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestClient_DeleteManagedLogging_validation(t *testing.T) {
	err := testClient.DeleteManagedLogging(&DeleteManagedLoggingInput{
		ServiceID: "",
		Kind:      ManagedLoggingInstanceOutput,
	})
	if err != ErrMissingServiceID {
		t.Errorf("unexpected error: %s", err)
	}

	err = testClient.DeleteManagedLogging(&DeleteManagedLoggingInput{
		ServiceID: testServiceID,
	})
	if err != ErrMissingKind {
		t.Errorf("unexpected error: %s", err)
	}

	err = testClient.DeleteManagedLogging(&DeleteManagedLoggingInput{
		ServiceID: testServiceID,
		Kind:      999,
	})
	if err != ErrNotImplemented {
		t.Errorf("unexpected error: %s", err)
	}
}
