package thresholds

import (
	"context"
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v12/fastly"
)

const (
	testAction   string = "block"
	testInterval int    = 60
	testLimit    int    = 1
	testName     string = "Test Threshold"
	testSignal   string = "SQLI"
)

func TestClient_Thresholds(t *testing.T) {
	var err error

	var thresholdID string

	// Create a test threshold.
	var threshold *Threshold

	fastly.Record(t, "create_threshold", func(c *fastly.Client) {
		threshold, err = Create(context.TODO(), c, &CreateInput{
			Action:      fastly.ToPointer(testAction),
			Interval:    fastly.ToPointer(testInterval),
			Limit:       fastly.ToPointer(testLimit),
			Name:        fastly.ToPointer(testName),
			Signal:      fastly.ToPointer(testSignal),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if threshold.Action != testAction {
		t.Errorf("unexpected threshold action: got %q, expected %q", threshold.Action, testAction)
	}
	if threshold.Interval != testInterval {
		t.Errorf("unexpected threshold interval: got %q, expected %q", threshold.Interval, testInterval)
	}
	if threshold.Limit != testLimit {
		t.Errorf("unexpected threshold limit: got %q, expected %q", threshold.Limit, testLimit)
	}
	if threshold.Name != testName {
		t.Errorf("unexpected threshold name: got %q, expected %q", threshold.Name, testName)
	}
	if threshold.Signal != testSignal {
		t.Errorf("unexpected threshold signal: got %q, expected %q", threshold.Signal, testSignal)
	}
	if threshold.DontNotify {
		t.Errorf("unexpected threshold don't notify: got %t, should default to false", threshold.DontNotify)
	}
	if threshold.Duration != 0 {
		t.Errorf("unexpected threshold duration: got %q, should default to 0", threshold.Duration)
	}
	if threshold.Enabled {
		t.Errorf("unexpected threshold enabled: got %t, should default to false", threshold.Enabled)
	}
	thresholdID = threshold.ThresholdID

	// Ensure we delete the test threshold at the end.
	defer func() {
		fastly.Record(t, "delete_threshold", func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				ThresholdID: fastly.ToPointer(thresholdID),
				WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			})
		})
		if err != nil {
			t.Errorf("error during threshold cleanup: %v", err)
		}
	}()

	// Get the test threshold.
	var getTestThreshold *Threshold
	fastly.Record(t, "get_threshold", func(c *fastly.Client) {
		getTestThreshold, err = Get(context.TODO(), c, &GetInput{
			ThresholdID: fastly.ToPointer(thresholdID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if getTestThreshold.Action != testAction {
		t.Errorf("unexpected threshold action: got %q, expected %q", getTestThreshold.Action, testAction)
	}
	if getTestThreshold.Interval != testInterval {
		t.Errorf("unexpected threshold interval: got %q, expected %q", getTestThreshold.Interval, testInterval)
	}
	if getTestThreshold.Limit != testLimit {
		t.Errorf("unexpected threshold limit: got %q, expected %q", getTestThreshold.Limit, testLimit)
	}
	if getTestThreshold.Name != testName {
		t.Errorf("unexpected threshold name: got %q, expected %q", getTestThreshold.Name, testName)
	}
	if getTestThreshold.Signal != testSignal {
		t.Errorf("unexpected threshold signal: got %q, expected %q", getTestThreshold.Signal, testSignal)
	}
	if getTestThreshold.DontNotify {
		t.Errorf("unexpected threshold don't notify: got %t, should default to false", getTestThreshold.DontNotify)
	}
	if getTestThreshold.Duration != 0 {
		t.Errorf("unexpected threshold duration: got %q, should default to 0", getTestThreshold.Duration)
	}
	if getTestThreshold.Enabled {
		t.Errorf("unexpected threshold enabled: got %t, should default to false", getTestThreshold.Enabled)
	}

	// Update the test threshold.
	const (
		updatedtestAction     string = "log"
		updatedTestDontNotify bool   = true
		updatedTestDuration   int    = 1
		updatedTestEnabled    bool   = true
		updatedTestInterval   int    = 600
		updateTestLimit       int    = 2
		updatedTestName       string = "Updated Test Threshold"
		updatedTestSignal     string = "BHH"
	)

	var updatedThreshold *Threshold
	fastly.Record(t, "update_threshold", func(c *fastly.Client) {
		updatedThreshold, err = Update(context.TODO(), c, &UpdateInput{
			Action:      fastly.ToPointer(updatedtestAction),
			Duration:    fastly.ToPointer(updatedTestDuration),
			DontNotify:  fastly.ToPointer(updatedTestDontNotify),
			Enabled:     fastly.ToPointer(updatedTestEnabled),
			Interval:    fastly.ToPointer(updatedTestInterval),
			Limit:       fastly.ToPointer(updateTestLimit),
			Name:        fastly.ToPointer(updatedTestName),
			Signal:      fastly.ToPointer(updatedTestSignal),
			ThresholdID: fastly.ToPointer(thresholdID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedThreshold.Action != updatedtestAction {
		t.Errorf("unexpected threshold action: got %q, expected %q", updatedThreshold.Action, updatedtestAction)
	}
	if updatedThreshold.Interval != updatedTestInterval {
		t.Errorf("unexpected threshold interval: got %q, expected %q", updatedThreshold.Interval, updatedTestInterval)
	}
	if updatedThreshold.Limit != updateTestLimit {
		t.Errorf("unexpected threshold limit: got %q, expected %q", updatedThreshold.Limit, updateTestLimit)
	}
	if updatedThreshold.Name != updatedTestName {
		t.Errorf("unexpected threshold name: got %q, expected %q", updatedThreshold.Name, updatedTestName)
	}
	if updatedThreshold.Signal != updatedTestSignal {
		t.Errorf("unexpected threshold signal: got %q, expected %q", updatedThreshold.Signal, updatedTestSignal)
	}
	if updatedThreshold.DontNotify != updatedTestDontNotify {
		t.Errorf("unexpected threshold don't notify: got %t, expected %t", updatedThreshold.DontNotify, updatedTestDontNotify)
	}
	if updatedThreshold.Duration != updatedTestDuration {
		t.Errorf("unexpected threshold duration: got %q, expected %q", updatedThreshold.Duration, updatedTestDuration)
	}
	if updatedThreshold.Enabled != updatedTestEnabled {
		t.Errorf("unexpected threshold enabled: got %t, expected %t", updatedThreshold.Enabled, updatedTestEnabled)
	}

	// List the thresholds for the test workspace and check the updated one is the only entry
	var listedThresholds *Thresholds
	fastly.Record(t, "list_threshold", func(c *fastly.Client) {
		listedThresholds, err = List(context.TODO(), c, &ListInput{
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			Limit:       fastly.ToPointer(1),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(listedThresholds.Data) != 1 {
		t.Errorf("unexpected threshold list length: got %q, expected %q", len(listedThresholds.Data), 1)
	}
	if listedThresholds.Data[0].ThresholdID != thresholdID {
		t.Errorf("unexpected threshold id: got %q, expected %q", listedThresholds.Data[0].ThresholdID, thresholdID)
	}
	if listedThresholds.Data[0].Signal != updatedTestSignal {
		t.Errorf("unexpected threshold field: got %q, expected %q", listedThresholds.Data[0].Signal, updatedTestSignal)
	}
}

func TestClient_GetThreshold_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		ThresholdID: nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingThresholdID) {
		t.Errorf("expected ErrMissingThresholdID: got %s", err)
	}
}

func TestClient_CreateThreshold_validation(t *testing.T) {
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
		t.Errorf("expected ErrMissingName: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Action:      nil,
		Name:        fastly.ToPointer(testName),
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingAction) {
		t.Errorf("expected ErrMissingAction: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Action:      fastly.ToPointer(testAction),
		Limit:       nil,
		Name:        fastly.ToPointer(testName),
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingLimit) {
		t.Errorf("expected ErrMissingLimit: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Action:      fastly.ToPointer(testAction),
		Limit:       fastly.ToPointer(testLimit),
		Interval:    nil,
		Name:        fastly.ToPointer(testName),
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingInterval) {
		t.Errorf("expected ErrMissingInterval: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Action:      fastly.ToPointer(testAction),
		Limit:       fastly.ToPointer(testLimit),
		Interval:    fastly.ToPointer(testInterval),
		Signal:      nil,
		Name:        fastly.ToPointer(testName),
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingSignal) {
		t.Errorf("expected ErrMissingSignal: got %s", err)
	}
}

func TestClient_UpdateThreshold_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		ThresholdID: nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingThresholdID) {
		t.Errorf("expected ErrMissingThresholdID: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		ThresholdID: fastly.ToPointer("someID"),
		Action:      nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingAction) {
		t.Errorf("expected ErrMissingAction: got %s", err)
	}
}

func TestClient_DeleteThreshold_validation(t *testing.T) {
	var err error
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		ThresholdID: nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingThresholdID) {
		t.Errorf("expected ErrMissingThresholdID: got %s", err)
	}
}
