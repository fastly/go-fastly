package timeseries

import (
	"errors"
	"testing"
	"time"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Global value for the tests.
var testWorkspaceID = fastly.TestNGWAFWorkspaceID

const tsMetrics = "XSS,SQLI,HTTP404"

var (
	// NOTE: Update this to a recent timestamp when regenerating the test fixtures,
	// otherwise the data may be outside of retention and an error will be
	// returned.
	tsEnd   = time.Date(2025, 6, 6, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	tsStart = time.Date(2025, 6, 4, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	tsGranularity = 60
)

func TestTime_Series(t *testing.T) {
	var err error
	var ts *TimeSeries

	// Get request timeseries metrics for given workspace
	fastly.Record(t, "get_timeseries", func(c *fastly.Client) {
		ts, err = Get(c, &GetInput{
			End:         &tsEnd,
			Granularity: &tsGranularity,
			Start:       &tsStart,
			Metrics:     fastly.ToPointer(tsMetrics),
			WorkspaceID: &testWorkspaceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ts == nil {
		t.Fatal("expected timeseries response, got nil")
	}
}

func TestClient_GetVirtualPatch_validation(t *testing.T) {
	var err error
	_, err = Get(fastly.TestClient, &GetInput{
		Start:       nil,
		Metrics:     fastly.ToPointer(tsMetrics),
		WorkspaceID: &testWorkspaceID,
	})
	if !errors.Is(err, fastly.ErrMissingStart) {
		t.Errorf("expected ErrMissingStart: got %s", err)
	}

	_, err = Get(fastly.TestClient, &GetInput{
		Start:       &tsStart,
		Metrics:     nil,
		WorkspaceID: &testWorkspaceID,
	})
	if !errors.Is(err, fastly.ErrMissingMetrics) {
		t.Errorf("expected ErrMissingMetrics: got %s", err)

		_, err = Get(fastly.TestClient, &GetInput{
			Start:       &tsStart,
			Metrics:     fastly.ToPointer(tsMetrics),
			WorkspaceID: nil,
		})
		if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
			t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
		}
	}
}
