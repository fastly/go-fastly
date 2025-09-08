package fastly

import (
	"context"
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
			_ = c.DeleteService(context.TODO(), &DeleteServiceInput{
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
		_, err = c.CreateDomain(context.TODO(), &CreateDomainInput{
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
		_, err = c.CreateDomain(context.TODO(), &CreateDomainInput{
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
		_, err = c.ActivateVersion(context.TODO(), &ActivateVersionInput{
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
		subscription, err = c.CreateTLSSubscription(context.TODO(), &CreateTLSSubscriptionInput{
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
			_ = c.DeleteTLSSubscription(context.TODO(), &DeleteTLSSubscriptionInput{
				ID: subscription.ID,
			})
		})
	}()

	// List
	var listSubscriptions []*TLSSubscription
	Record(t, fixtureBase+"list", func(c *Client) {
		listSubscriptions, err = c.ListTLSSubscriptions(context.TODO(), &ListTLSSubscriptionsInput{
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
		retrievedSubscription, err = c.GetTLSSubscription(context.TODO(), &GetTLSSubscriptionInput{
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
		updatedSubscription, err = c.UpdateTLSSubscription(context.TODO(), &UpdateTLSSubscriptionInput{
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
		err = c.DeleteTLSSubscription(context.TODO(), &DeleteTLSSubscriptionInput{
			ID: subscription.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_TLSSubscription_Compute(t *testing.T) {
	t.Parallel()

	// NOTE: TLS Subscriptions require the domains specified to exist on an
	// activated service version.

	var err error
	var tv *Version
	Record(t, fixtureBase+"compute/version", func(c *Client) {
		tv = testVersionCompute(t, c)
	})

	// Ensure service (and all domains within it) are deleted
	defer func() {
		Record(t, fixtureBase+"compute/version", func(c *Client) {
			_ = c.DeleteService(context.TODO(), &DeleteServiceInput{
				ServiceID: TestComputeServiceID,
			})
		})
	}()

	// Create domains needed to support the TLS Subscription tests.
	//
	// NOTE: We don't reuse the fixtures from the domains test file as we don't
	// want to create a complex coupling that could cause confusion in the future
	// if the domains in either test file were to change.

	domain1 := "integ-test3.go-fastly-3.com"
	domain2 := "integ-test4.go-fastly-4.com"

	Record(t, fixtureBase+"compute/domains/create", func(c *Client) {
		_, err = c.CreateDomain(context.TODO(), &CreateDomainInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer(domain1),
			Comment:        ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, fixtureBase+"compute/domains/create2", func(c *Client) {
		_, err = c.CreateDomain(context.TODO(), &CreateDomainInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer(domain2),
			Comment:        ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// a valid package must exist for compute services for activation
	// Update with valid package file path
	RecordIgnoreBody(t, fixtureBase+"compute/package_update", func(c *Client) {
		_, err = c.UpdatePackage(context.TODO(), &UpdatePackageInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			PackagePath:    ToPointer("test_assets/package/valid.tar.gz"),
		})
	})

	// Create backend
	Record(t, fixtureBase+"compute/backend/create", func(c *Client) {
		_, err = c.CreateBackend(context.TODO(), &CreateBackendInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-backend-compute"),
			Address:        ToPointer("integ-test3.go-fastly-3.com"),
			ConnectTimeout: ToPointer(1500),
			OverrideHost:   ToPointer("origin.example.com"),
			SSLCheckCert:   ToPointer(Compatibool(false)),
			SSLCiphers:     ToPointer("DHE-RSA-AES256-SHA:DHE-RSA-CAMELLIA256-SHA:AES256-GCM-SHA384"),
			SSLSNIHostname: ToPointer("ssl-hostname.com"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Activate service version otherwise TLS Subscription won't be able to locate	// the specified domains and the API will return an error.

	Record(t, fixtureBase+"compute/activate_version", func(c *Client) {
		_, err = c.ActivateVersion(context.TODO(), &ActivateVersionInput{
			ServiceID:      TestComputeServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create
	var subscription *TLSSubscription
	Record(t, fixtureBase+"compute/create", func(c *Client) {
		subscription, err = c.CreateTLSSubscription(context.TODO(), &CreateTLSSubscriptionInput{
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
		Record(t, fixtureBase+"compute/cleanup", func(c *Client) {
			// NOTE: We would expect this API call to produce a 404 rather than a 204
			// because the "delete" step at the end of the test function is
			// effectively deleting the subscription, and then this defer function is
			// executing. Meaning we've already deleted the subscription. The defer
			// function is here to ensure the subscription is deleted, just in case
			// any of the other API calls unexpectedly fail before the "delete" step
			// at the end of the test.
			_ = c.DeleteTLSSubscription(context.TODO(), &DeleteTLSSubscriptionInput{
				ID: subscription.ID,
			})

			_ = c.DeleteBackend(context.TODO(), &DeleteBackendInput{
				ServiceID:      TestComputeServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-backend-compute",
			})
		})
	}()

	// List
	var listSubscriptions []*TLSSubscription
	Record(t, fixtureBase+"compute/list", func(c *Client) {
		listSubscriptions, err = c.ListTLSSubscriptions(context.TODO(), &ListTLSSubscriptionsInput{
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
	Record(t, fixtureBase+"compute/get", func(c *Client) {
		retrievedSubscription, err = c.GetTLSSubscription(context.TODO(), &GetTLSSubscriptionInput{
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
	Record(t, fixtureBase+"compute/update", func(c *Client) {
		updatedSubscription, err = c.UpdateTLSSubscription(context.TODO(), &UpdateTLSSubscriptionInput{
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

	Record(t, fixtureBase+"compute/delete", func(c *Client) {
		err = c.DeleteTLSSubscription(context.TODO(), &DeleteTLSSubscriptionInput{
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
		tlsSubscriptions, err = c.ListTLSSubscriptions(context.TODO(), &ListTLSSubscriptionsInput{})
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
		_, err = c.CreateTLSSubscription(context.TODO(), &CreateTLSSubscriptionInput{
			Domains: []*TLSDomain{
				{ID: "DOMAIN_NAME"},
			},
			CommonName: &TLSDomain{ID: "DOMAIN_NAME"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = TestClient.CreateTLSSubscription(context.TODO(), &CreateTLSSubscriptionInput{})
	if !errors.Is(err, ErrMissingTLSDomain) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateTLSSubscription(context.TODO(), &CreateTLSSubscriptionInput{
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

	_, err := TestClient.GetTLSSubscription(context.TODO(), &GetTLSSubscriptionInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateTLSSubscription_validation(t *testing.T) {
	t.Parallel()

	_, err := TestClient.UpdateTLSSubscription(context.TODO(), &UpdateTLSSubscriptionInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteTLSSubscription_validation(t *testing.T) {
	t.Parallel()

	err := TestClient.DeleteTLSSubscription(context.TODO(), &DeleteTLSSubscriptionInput{})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}
