package signals

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/common"
)

const (
	testName        = "A Real Name "
	testDescription = "This is a description"
	testRefID       = "a-real-name-"
)

func TestClient_Signals_WorkspaceScope(t *testing.T) {
	runSignalsTest(t, common.ScopeTypeWorkspace, fastly.TestNGWAFWorkspaceID)
}

func TestClient_Signals_AccountScope(t *testing.T) {
	runSignalsTest(t, common.ScopeTypeAccount, "*") // assuming TestNGWAFAccountID exists
}

func runSignalsTest(t *testing.T, scopeType common.ScopeType, appliesToID string) {
	var err error
	var signalID string
	testSignalName := testName + string(scopeType)
	testSignalReferenceID := testRefID + string(scopeType)

	// Create a test signal.
	var signal *Signal

	fastly.Record(t, fmt.Sprintf("%s_create_signal", scopeType), func(c *fastly.Client) {
		signal, err = Create(context.TODO(), c, &CreateInput{
			Description: fastly.ToPointer(testDescription),
			Name:        fastly.ToPointer(testSignalName),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if signal.Name != testSignalName {
		t.Errorf("unexpected signal name: got %q, expected %q", signal.Name, testSignalName)
	}
	if signal.Description != testDescription {
		t.Errorf("unexpected signal description: got %q, expected %q", signal.Description, testDescription)
	}
	if !strings.Contains(signal.ReferenceID, testSignalReferenceID) {
		t.Errorf("unexpected list reference ID: expected %q to contain %q", signal.ReferenceID, testSignalReferenceID)
	}
	signalID = signal.SignalID

	// Ensure we delete the test signal at the end.
	defer func() {
		fastly.Record(t, fmt.Sprintf("%s_delete_signal", scopeType), func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				Scope: &common.Scope{
					Type:      scopeType,
					AppliesTo: []string{appliesToID},
				},
				SignalID: fastly.ToPointer(signalID),
			})
		})
		if err != nil {
			t.Errorf("error during signal cleanup: %v", err)
		}
	}()

	// Get the test signal.
	var getTestSignal *Signal
	fastly.Record(t, fmt.Sprintf("%s_get_signal", scopeType), func(c *fastly.Client) {
		getTestSignal, err = Get(context.TODO(), c, &GetInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			SignalID: fastly.ToPointer(signalID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if getTestSignal.Name != testSignalName {
		t.Errorf("unexpected signal name: got %q, expected %q", getTestSignal.Name, testSignalName)
	}
	if getTestSignal.Description != testDescription {
		t.Errorf("unexpected signal description: got %q, expected %q", getTestSignal.Description, testDescription)
	}

	// Update the test signal.
	const updatedSignalDescription = "This is an updated description"

	var updatedSignal *Signal
	fastly.Record(t, fmt.Sprintf("%s_update_signal", scopeType), func(c *fastly.Client) {
		updatedSignal, err = Update(context.TODO(), c, &UpdateInput{
			Description: fastly.ToPointer(string(updatedSignalDescription)),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			SignalID: fastly.ToPointer(signalID),
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
	fastly.Record(t, fmt.Sprintf("%s_list_signals", scopeType), func(c *fastly.Client) {
		listedSignals, err = List(context.TODO(), c, &ListInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
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
	if listedSignals.Data[0].Name != testSignalName {
		t.Errorf("unexpected signal field: got %q, expected %q", listedSignals.Data[0].Name, testSignalName)
	}
	if !strings.Contains(signal.ReferenceID, testSignalReferenceID) {
		t.Errorf("unexpected list reference ID: expected %q to contain %q", signal.ReferenceID, testSignalReferenceID)
	}
	if listedSignals.Data[0].Scope.Type != string(scopeType) {
		t.Errorf("unexpected signal scope type: got %q, expected %q", listedSignals.Data[0].Scope.Type, string(scopeType))
	}
	if listedSignals.Data[0].Scope.AppliesTo[0] != appliesToID {
		t.Errorf("unexpected signal scope appliesTo: got %q, expected %q", listedSignals.Data[0].Scope.AppliesTo[0], appliesToID)
	}
}

func TestClient_GetSignal_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		Scope: nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		SignalID: nil,
		Scope: &common.Scope{
			Type:      common.ScopeTypeWorkspace,
			AppliesTo: []string{},
		},
	})
	if !errors.Is(err, fastly.ErrMissingSignalID) {
		t.Errorf("expected ErrMissingSignalID: got %s", err)
	}
}

func TestClient_CreateSignal_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Name: nil,
	})
	if !errors.Is(err, fastly.ErrMissingName) {
		t.Errorf("expected ErrMissingName: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Name:  fastly.ToPointer("some Name"),
		Scope: nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}

func TestClient_UpdateSignal_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		SignalID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingSignalID) {
		t.Errorf("expected ErrMissingSignalID: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		SignalID: fastly.ToPointer("someID"),
		Scope:    nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		Description: nil,
		SignalID:    fastly.ToPointer("someID"),
		Scope: &common.Scope{
			Type:      common.ScopeTypeWorkspace,
			AppliesTo: []string{},
		},
	})
	if !errors.Is(err, fastly.ErrMissingDescription) {
		t.Errorf("expected ErrMissingField: got %s", err)
	}
}

func TestClient_DeleteSignal_validation(t *testing.T) {
	var err error
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		SignalID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingSignalID) {
		t.Errorf("expected ErrMissingSignalID: got %s", err)
	}
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		SignalID: fastly.ToPointer("someID"),
		Scope:    nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}
