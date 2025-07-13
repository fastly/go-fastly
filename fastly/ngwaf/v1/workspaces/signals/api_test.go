package signals

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v11/fastly"
)

const (
	testName        = "A real Name"
	testDescription = "This is a description"
	testRefID       = "site.a-real-name"
)

func TestClient_Signals(t *testing.T) {
	var err error
	var signalID string

	// Create a test signal.
	var signal *Signal

	fastly.Record(t, "create_signal", func(c *fastly.Client) {
		signal, err = Create(context.TODO(), c, &CreateInput{
			Description: fastly.ToPointer(testDescription),
			Name:        fastly.ToPointer(testName),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if signal.Name != testName {
		t.Errorf("unexpected signal name: got %q, expected %q", signal.Name, testName)
	}
	if signal.Description != testDescription {
		t.Errorf("unexpected signal description: got %q, expected %q", signal.Description, testDescription)
	}
	if signal.ReferenceID != testRefID {
		t.Errorf("unexpected signal referenceID: got %q, expected %q", signal.ReferenceID, testRefID)
	}
	signalID = signal.SignalID

	// Ensure we delete the test signal at the end.
	defer func() {
		fastly.Record(t, "delete_signal", func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				SignalID:    fastly.ToPointer(signalID),
				WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			})
		})
		if err != nil {
			t.Errorf("error during signal cleanup: %v", err)
		}
	}()

	// Get the test signal.
	var getTestSignal *Signal
	fastly.Record(t, "get_signal", func(c *fastly.Client) {
		getTestSignal, err = Get(context.TODO(), c, &GetInput{
			SignalID:    fastly.ToPointer(signalID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if getTestSignal.Name != testName {
		t.Errorf("unexpected signal name: got %q, expected %q", getTestSignal.Name, testName)
	}
	if getTestSignal.Description != testDescription {
		t.Errorf("unexpected signal description: got %q, expected %q", getTestSignal.Description, testDescription)
	}

	// Update the test signal.
	const updatedSignalDescription = "This is an updated description"

	var updatedSignal *Signal
	fastly.Record(t, "update_signal", func(c *fastly.Client) {
		updatedSignal, err = Update(context.TODO(), c, &UpdateInput{
			Description: fastly.ToPointer(string(updatedSignalDescription)),
			SignalID:    fastly.ToPointer(signalID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedSignal.Description != updatedSignalDescription {
		t.Errorf("unexpected signal type: got %q, expected %q", updatedSignal.Description, updatedSignalDescription)
	}

	// List the signals for the test workspace and check the updated one is the only entry
	var listedSignals *Signals
	fastly.Record(t, "list_signal", func(c *fastly.Client) {
		listedSignals, err = List(context.TODO(), c, &ListInput{
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(listedSignals.Data) != 1 {
		t.Errorf("unexpected signal list length: got %q, expected %q", len(listedSignals.Data), 1)
	}
	if listedSignals.Data[0].Description != updatedSignalDescription {
		t.Errorf("unexpected signal type: got %q, expected %q", listedSignals.Data[0].Description, updatedSignalDescription)
	}
	if listedSignals.Data[0].Name != testName {
		t.Errorf("unexpected signal field: got %q, expected %q", listedSignals.Data[0].Name, testName)
	}
}

func TestClient_GetSignal_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		SignalID:    nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingSignalID) {
		t.Errorf("expected ErrMissingSignalID: got %s", err)
	}
}

func TestClient_CreateSignal_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Name:        nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingName) {
		t.Errorf("expected ErrMissingField: got %s", err)
	}
}

func TestClient_UpdateSignal_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		SignalID:    nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingSignalID) {
		t.Errorf("expected ErrMissingSignalID: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		Description: nil,
		SignalID:    fastly.ToPointer("someID"),
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingDescription) {
		t.Errorf("expected ErrMissingField: got %s", err)
	}
}

func TestClient_DeleteSignal_validation(t *testing.T) {
	var err error
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		SignalID:    nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingSignalID) {
		t.Errorf("expected ErrMissingSignalID: got %s", err)
	}
}
