package status

import (
	"errors"
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v10/fastly"
)

func TestClient_DomainToolsStatus(t *testing.T) {
	t.Parallel()

	var err error
	var status *Status
	domain := "fastly-sdk-gofastly-testing.com"
	fastly.Record(t, "get", func(client *fastly.Client) {
		status, err = Get(client, &GetInput{
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

	if len(status.Offers) > 0 {
		t.Errorf("incorrect offers, expected %d, got %d", 0, len(status.Offers))
	}

	if !strings.Contains(status.Status, "inactive") {
		t.Errorf("incorrect status, expected %s within status, got %s", "inactive", status.Status)
	}

	if !strings.Contains(status.Tags, "generic") {
		t.Errorf("incorrect tags, expected %s within tags, got %s", "generic", status.Tags)
	}

	if status.Scope != nil {
		t.Errorf("incorrect scope, unexpected presence of scope, got %s", *status.Scope)
	}
}

func TestClient_DomainToolsStatusEstimate(t *testing.T) {
	t.Parallel()

	var err error
	var status *Status
	domain := "fastly-sdk-gofastly-testing.com"
	fastly.Record(t, "get_estimate", func(client *fastly.Client) {
		status, err = Get(client, &GetInput{
			Domain: domain,
			Scope:  fastly.ToPointer(ScopeEstimate),
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

	if status.Scope == nil || *status.Scope != ScopeEstimate {
		t.Errorf("incorrect scope, expected %s, got %v", ScopeEstimate, status.Scope)
	}
}

func TestClient_DomainToolsStatusOffers(t *testing.T) {
	t.Parallel()

	var err error
	var status *Status
	domain := "domainr-testing.org"
	fastly.Record(t, "get_offers", func(client *fastly.Client) {
		status, err = Get(client, &GetInput{
			Domain: domain,
			Scope:  fastly.ToPointer(ScopeEstimate),
		})
	})

	if err != nil {
		t.Fatal(err)
	}

	if status.Domain != domain {
		t.Errorf("incorrect domain, expected %s, got %s", domain, status.Domain)
	}

	if status.Zone != "org" {
		t.Errorf("incorrect zone, expected %s, got %s", "org", status.Zone)
	}

	if status.Scope == nil || *status.Scope != ScopeEstimate {
		t.Errorf("incorrect scope, expected %s, got %v", ScopeEstimate, status.Scope)
	}

	if !strings.Contains(status.Status, "priced") {
		t.Errorf("incorrect status, expected %s within status, got %s", "priced", status.Status)
	}

	if len(status.Offers) == 0 {
		t.Errorf("incorrect offers, expected at least one offer, got %d", len(status.Offers))
	}

	if status.Offers[0].Currency != "USD" {
		t.Errorf("incorrect currency, expected %s, got %s", "USD", status.Offers[0].Currency)
	}
}

func TestClient_DomainToolsStatusErrorHandling(t *testing.T) {
	t.Parallel()

	var err error
	domain := "fastly-sdk-gofastly-testing"
	fastly.Record(t, "get_error", func(client *fastly.Client) {
		_, err = Get(client, &GetInput{
			Domain: domain,
		})
	})

	if err == nil {
		t.Errorf("no error returned")
	}

	httpErr := &fastly.HTTPError{}
	if !errors.As(err, &httpErr) {
		t.Errorf("wrong error returned")
	}

	if len(httpErr.Errors) != 1 {
		t.Errorf("wrong number of errors returned")
	}

	if httpErr.Errors[0].Status != "400" {
		t.Errorf("bad error status, expected %s, got %s", "400", httpErr.Errors[0].Status)
	}

	if httpErr.Errors[0].Title != "Domain not found" {
		t.Errorf("bad error message, expected %s, got %s", "domain not found", httpErr.Errors[0].Title)
	}

	if httpErr.Errors[0].Detail != "fastly-sdk-gofastly-testing" {
		t.Errorf("bad detail, expected %s, got %s", "fastly-sdk-gofastly-testing", httpErr.Errors[0].Detail)
	}
}
