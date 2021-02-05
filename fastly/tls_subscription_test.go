package fastly

import "testing"

const fixtureBase = "tls_subscription/"

func TestClient_TLSSubscription(t *testing.T) {
	t.Parallel()

	// Create
	var err error
	var subscription *TLSSubscription
	record(t, fixtureBase+"create", func(c *Client) {
		subscription, err = c.CreateTLSSubscription(&CreateTLSSubscriptionInput{
			Domains: []*TLSDomain{
				{ID: "DOMAIN_NAME"},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteTLSSubscription(&DeleteTLSSubscriptionInput{
				ID: subscription.ID,
			})
		})
	}()

	// List
	var listSubscriptions []*TLSSubscription
	record(t, fixtureBase+"list", func(c *Client) {
		listSubscriptions, err = c.ListTLSSubscriptions(&ListTLSSubscriptionsInput{})
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
	record(t, fixtureBase+"get", func(c *Client) {
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

	record(t, fixtureBase+"delete", func(c *Client) {
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
	record(t, fixtureBase+"list", func(c *Client) {
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
	record(t, fixtureBase+"create", func(c *Client) {
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

	_, err = testClient.CreateTLSSubscription(&CreateTLSSubscriptionInput{})
	if err != ErrMissingTLSDomain {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateTLSSubscription(&CreateTLSSubscriptionInput{
		Domains: []*TLSDomain{
			{ID: "DN1"},
			{ID: "DN2"},
		},
		CommonName: &TLSDomain{ID: "DN3"},
	})
	if err != ErrCommonNameNotInDomains {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetTLSSubscription_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, fixtureBase+"get", func(c *Client) {
		_, err = c.GetTLSSubscription(&GetTLSSubscriptionInput{
			ID: "SUBSCRIPTION_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.GetTLSSubscription(&GetTLSSubscriptionInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteTLSSubscription_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteTLSSubscription(&DeleteTLSSubscriptionInput{
			ID: "SUBSCRIPTION_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	err = testClient.DeleteTLSSubscription(&DeleteTLSSubscriptionInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
