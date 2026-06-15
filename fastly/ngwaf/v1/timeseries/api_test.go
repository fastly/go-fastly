package timeseries

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fastly/go-fastly/v15/fastly"
)

const tsMetrics = "XSS,SQLI,HTTP404"

var (
	// NOTE: Update this to a recent timestamp when regenerating the test fixtures,
	// otherwise the data may be outside of retention and an error will be
	// returned.
	tsTo   = time.Date(2026, 6, 15, 12, 10, 0, 0, time.UTC).Format(time.RFC3339)
	tsFrom = time.Date(2026, 6, 15, 11, 58, 0, 0, time.UTC).Format(time.RFC3339)

	// Set to 60 seconds of granularity.
	tsGranularity = 60
)

func TestTimeSeries_List(t *testing.T) {
	var err error
	var ts *Timeseries

	fastly.Record(t, "list_timeseries", func(c *fastly.Client) {
		ts, err = List(context.TODO(), c, &ListInput{
			From:        &tsFrom,
			Granularity: &tsGranularity,
			Metrics:     fastly.ToPointer(tsMetrics),
			To:          &tsTo,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ts == nil {
		t.Fatal("expected timeseries response, got nil")
	}
}

func TestClient_List_Timeseries_validation(t *testing.T) {
	var err error
	_, err = List(context.TODO(), fastly.TestClient, &ListInput{
		From:    nil,
		Metrics: fastly.ToPointer(tsMetrics),
	})
	if !errors.Is(err, fastly.ErrMissingFrom) {
		t.Errorf("expected ErrMissingFrom: got %s", err)
	}

	_, err = List(context.TODO(), fastly.TestClient, &ListInput{
		From:    &tsFrom,
		Metrics: nil,
	})
	if !errors.Is(err, fastly.ErrMissingMetrics) {
		t.Errorf("expected ErrMissingMetrics: got %s", err)
	}
}
