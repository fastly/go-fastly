package status

import (
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

func TestClient_DomainToolsStatus(t *testing.T) {
	t.Parallel()

	var err error
	var status *Status
	domain := "fastly-sdk-gofastly-testing.com"
	fastly.Record(t, "get", func(client *fastly.Client) {
		status, err = Get(client, &Input{
			Domain: domain,
		})
	})

	if err != nil {
		t.Fatal(err)
	}

	if status.Domain != domain {
		t.Errorf("incorrect domain, expected %s, got %s", domain, status.Domain)
	}

	if status.Zone != "com" {
		t.Errorf("incorrect zone, expected %s, got %s", "com", status.Zone)
	}
}
