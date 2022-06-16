package fastly

import (
	"testing"
)

func TestClient_ServiceAuthorizations(t *testing.T) {
	t.Parallel()

	fixtureBase := "service_authorizations/"

	// Create
	var err error
	var sa *ServiceAuthorization
	record(t, fixtureBase+"create", func(c *Client) {
		sa, err = c.CreateServiceAuthorization(&CreateServiceAuthorizationInput{
			ServiceID:  "7ZVxm5pPWdzKdl3P5UW7jR",
			UserID:     "4tKBSuFhNEiIpNDxmmVydt",
			Permission: "full",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteServiceAuthorization(&DeleteServiceAuthorizationInput{
				ID: sa.ID,
			})

		})
	}()

	if sa.ServiceID != "7ZVxm5pPWdzKdl3P5UW7jR" {
		t.Errorf("bad service id: %v", sa.ServiceID)
	}

	if sa.UserID != "4tKBSuFhNEiIpNDxmmVydt" {
		t.Errorf("bad user id: %v", sa.UserID)
	}

	if sa.Permission != "full" {
		t.Errorf("bad permission: %v", sa.Permission)
	}

	// Get
	var nsa *ServiceAuthorization
	record(t, fixtureBase+"get", func(c *Client) {
		nsa, err = c.GetServiceAuthorization(&GetServiceAuthorizationInput{
			ID: sa.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if nsa.ServiceID != "7ZVxm5pPWdzKdl3P5UW7jR" {
		t.Errorf("bad service id: %v", nsa.ServiceID)
	}

	// Update
	var usa *ServiceAuthorization
	record(t, fixtureBase+"update", func(c *Client) {
		usa, err = c.UpdateServiceAuthorization(&UpdateServiceAuthorizationInput{
			ID:          sa.ID,
			Permissions: "purge_select",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if usa.ServiceID != "7ZVxm5pPWdzKdl3P5UW7jR" {
		t.Errorf("bad service id: %v", usa.ServiceID)
	}
	if usa.Permission != "purge_select" {
		t.Errorf("bad permission: %v", usa.Permission)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteServiceAuthorization(&DeleteServiceAuthorizationInput{
			ID: sa.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetServiceAuthorization_validation(t *testing.T) {
	var err error
	_, err = testClient.GetServiceAuthorization(&GetServiceAuthorizationInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateServiceAuthorization_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateServiceAuthorization(&CreateServiceAuthorizationInput{
		ServiceID: "",
		UserID:    "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateServiceAuthorization(&CreateServiceAuthorizationInput{
		ServiceID: "my-service-id",
		UserID:    "",
	})
	if err != ErrMissingUserID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateServiceAuthorization_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateServiceAuthorization(&UpdateServiceAuthorizationInput{
		ID:          "",
		Permissions: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateServiceAuthorization(&UpdateServiceAuthorizationInput{
		ID:          "my-service-authorization-id",
		Permissions: "",
	})
	if err != ErrMissingPermissions {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteServiceAuthorization_validation(t *testing.T) {
	err := testClient.DeleteServiceAuthorization(&DeleteServiceAuthorizationInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

}
