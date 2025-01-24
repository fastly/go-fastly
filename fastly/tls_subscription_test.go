package fastly

import (
	"errors"
	"testing"
)

const fixtureBase = "tls_subscription/"

func TestClient_TLSSubscription(t *testing.T) {
	t.Parallel()

	// NOTE: TLS Subscriptions require the domains specified to exist on an
	// activated service version.

	var err error
	var tv *Version
	Record(t, fixtureBase+"version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Ensure service (and all domains within it) are deleted
	defer func() {
		Record(t, fixtureBase+"version", func(c *Client) {
			_ = c.DeleteService(&DeleteServiceInput{
				ServiceID: TestDeliveryServiceID,
			})
		})
	}()

	// Create domains needed to support the TLS Subscription tests.
	//
	// NOTE: We don't reuse the fixtures from the domains test file as we don't
	// want to create a complex coupling that could cause confusion in the future
	// if the domains in either test file were to change.

	domain1 := "integ-test1.go-fastly-1.com"
	domain2 := "integ-test2.go-fastly-2.com"

	Record(t, fixtureBase+"domains/create", func(c *Client) {
		_, err = c.CreateDomain(&CreateDomainInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer(domain1),
			Comment:        ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, fixtureBase+"domains/create2", func(c *Client) {
		_, err = c.CreateDomain(&CreateDomainInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer(domain2),
			Comment:        ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Activate service version otherwise TLS Subscription won't be able to locate
	// the specified domains and the API will return an error.

	Record(t, fixtureBase+"activate_version", func(c *Client) {
		_, err = c.ActivateVersion(&ActivateVersionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create
	var subscription *TLSSubscription
	Record(t, fixtureBase+"create", func(c *Client) {
		subscription, err = c.CreateTLSSubscription(&CreateTLSSubscriptionInput{
			Domains: []*TLSDomain{
				{ID: domain1},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, fixtureBase+"cleanup", func(c *Client) {
			// NOTE: We would expect this API call to produce a 404 rather than a 204
			// because the "delete" step at the end of the test function is
			// effectively deleting the subscription, and then this defer function is
			// executing. Meaning we've already deleted the subscription. The defer
			// function is here to ensure the subscription is deleted, just in case
			// any of the other API calls unexpectedly fail before the "delete" step
			// at the end of the test.
			_ = c.DeleteTLSSubscription(&DeleteTLSSubscriptionInput{
				ID: subscription.ID,
			})
		})
	}()

	// List
	var listSubscriptions []*TLSSubscription
	Record(t, fixtureBase+"list", func(c *Client) {
		listSubscriptions, err = c.ListTLSSubscriptions(&ListTLSSubscriptionsInput{
			// NOTE: Added this filter so I could manually verify that the filter is
			// only added to the API request query parameters when set to `true`. See
			// notes in formatFilters for input struct for details of a possible API
			// bug regarding this filter.
			FilterActiveOrders: false,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(listSubscriptions) < 1 {
		t.Errorf("bad TLS subscriptions: %v", listSubscriptions)
	}
	if listSubscriptions[0].Domains == nil {
		t.Errorf("TLS Domains relation should not be nil: %v", listSubscriptions)
	}
	if len(listSubscriptions[0].Domains) < 1 {
		t.Errorf("TLS Domains list should not be empty: %v", listSubscriptions)
	}
	if listSubscriptions[0].Domains[0].ID != subscription.Domains[0].ID {
		t.Errorf("bad Domain ID: %s (%s)", listSubscriptions[0].Domains[0].ID, subscription.Domains[0].ID)
	}
	if listSubscriptions[0].CommonName.ID != subscription.Domains[0].ID {
		t.Errorf("bad CommonName: %s (%s)", listSubscriptions[0].CommonName.ID, subscription.Domains[0].ID)
	}

	var retrievedSubscription *TLSSubscription
	Record(t, fixtureBase+"get", func(c *Client) {
		retrievedSubscription, err = c.GetTLSSubscription(&GetTLSSubscriptionInput{
			ID: subscription.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if subscription.ID != retrievedSubscription.ID {
		t.Errorf("bad ID: %s (%s)", subscription.ID, retrievedSubscription.ID)
	}

	var updatedSubscription *TLSSubscription
	Record(t, fixtureBase+"update", func(c *Client) {
		updatedSubscription, err = c.UpdateTLSSubscription(&UpdateTLSSubscriptionInput{
			ID: subscription.ID,
			Domains: []*TLSDomain{
				{ID: domain1},
				{ID: domain2},
			},
			CommonName: &TLSDomain{ID: domain2},
			Force:      true,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedSubscription.CommonName.ID != domain2 {
		t.Errorf("bad CommonName %s (%s)", updatedSubscription.CommonName.ID, domain2)
	}

	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteTLSSubscription(&DeleteTLSSubscriptionInput{
			ID: subscription.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListTLSSubscriptions_validation(t *testing.T) {
	t.Parallel()

	var tlsSubscriptions []*TLSSubscription
	var err error
	Record(t, fixtureBase+"list", func(c *Client) {
		tlsSubscriptions, err = c.ListTLSSubscriptions(&ListTLSSubscriptionsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tlsSubscriptions) < 1 {
		t.Errorf("bad tls subscriptions: %v", tlsSubscriptions)
	}
}

func TestClient_CreateTLSSubscription_validation(t *testing.T) {
	t.Parallel()

	var err error
	Record(t, fixtureBase+"create", func(c *Client) {
		_, err = c.CreateTLSSubscription(&CreateTLSSubscriptionInput{
			Domains: []*TLSDomain{
				{ID: "DOMAIN_NAME"},
			},
			CommonName: &TLSDomain{ID: "DOMAIN_NAME"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = TestClient.CreateTLSSubscription(&CreateTLSSubscriptionInput{})
	if !errors.Is(err, ErrMissingTLSDomain) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateTLSSubscription(&CreateTLSSubscriptionInput{
		Domains: []*TLSDomain{
			{ID: "DN1"},
			{ID: "DN2"},
		},
		CommonName: &TLSDomain{ID: "DN3"},
	})
	if !errors.Is(err, ErrCommonNameNotInDomains) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetTLSSubscription_validation(t *testing.T) {
	t.Parallel()

	_, err := TestClient.GetTLSSubscription(&GetTLSSubscriptionInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateTLSSubscription_validation(t *testing.T) {
	t.Parallel()

	_, err := TestClient.UpdateTLSSubscription(&UpdateTLSSubscriptionInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteTLSSubscription_validation(t *testing.T) {
	t.Parallel()

	err := TestClient.DeleteTLSSubscription(&DeleteTLSSubscriptionInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}
