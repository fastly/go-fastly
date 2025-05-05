package suggest

import (
	"errors"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

func TestClient_Suggestion(t *testing.T) {
	t.Parallel()

	var err error
	var suggestions *Suggestions
	fastly.Record(t, "get", func(client *fastly.Client) {
		suggestions, err = Get(client, &Input{
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

	var comFound bool

	for i := 0; i < len(suggestions.Results); i++ {
		if suggestions.Results[i].Zone == "com" {
			comFound = true
		}
	}

	if !comFound {
		t.Errorf("no .com zone suggestion found in %d results", len(suggestions.Results))
	}

	fastly.Record(t, "get", func(client *fastly.Client) {
		suggestions, err = Get(client, &Input{})
	})

	if !errors.Is(err, fastly.ErrMissingDomainQuery) {
		t.Errorf("expected error but got %v", err)
	}
}
