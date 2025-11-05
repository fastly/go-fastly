package suggest

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/fastly/go-fastly/v12/fastly"
)

func TestClient_DomainToolsSuggestion(t *testing.T) {
	t.Parallel()

	var err error
	var suggestions *Suggestions
	fastly.Record(t, "get", func(client *fastly.Client) {
		suggestions, err = Get(context.TODO(), client, &GetInput{
			Query:    "fastly testing",
			Defaults: fastly.ToPointer("com"),
			Keywords: fastly.ToPointer("testing"),
		})
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(suggestions.Results) == 0 {
		t.Error("no suggestions found")
	}

	zoneFound := slices.ContainsFunc(suggestions.Results, func(s Suggestion) bool {
		return s.Zone == "com"
	})

	if !zoneFound {
		t.Errorf("no com zone suggestion found in %d suggestions", len(suggestions.Results))
	}

	// omit Query from GetInput
	fastly.Record(t, "get", func(client *fastly.Client) {
		suggestions, err = Get(context.TODO(), client, &GetInput{})
	})

	if !errors.Is(err, fastly.ErrMissingDomainQuery) {
		t.Errorf("expected error but got %v", err)
	}
}
