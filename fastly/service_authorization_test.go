package fastly

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/jsonapi"
)

func TestClient_ServiceAuthorizations(t *testing.T) {
	t.Parallel()

	fixtureBase := "service_authorizations/"

	// Create
	var err error
	var sa *ServiceAuthorization
	Record(t, fixtureBase+"create", func(c *Client) {
		sa, err = c.CreateServiceAuthorization(&CreateServiceAuthorizationInput{
			Service:    &SAService{ID: TestDeliveryServiceID},
			User:       &SAUser{ID: "1pnpEMCscfjqgvH7Qofda6"},
			Permission: "full",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// List
	var sasResp *ServiceAuthorizations
	Record(t, fixtureBase+"/list", func(c *Client) {
		sasResp, err = c.ListServiceAuthorizations(&ListServiceAuthorizationsInput{
			PageSize: 10,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sasResp.Items) == 0 {
		t.Errorf("bad service authorizations: %v", sasResp.Items)
	}

	// Ensure deleted
	defer func() {
		Record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteServiceAuthorization(&DeleteServiceAuthorizationInput{
				ID: sa.ID,
			})
		})
	}()

	if sa.Service.ID != TestDeliveryServiceID {
		t.Errorf("bad service id: %v", sa.Service.ID)
	}

	if sa.User.ID != "1pnpEMCscfjqgvH7Qofda6" {
		t.Errorf("bad user id: %v", sa.User.ID)
	}

	if sa.Permission != "full" {
		t.Errorf("bad permission: %v", sa.Permission)
	}

	// Get
	var nsa *ServiceAuthorization
	Record(t, fixtureBase+"get", func(c *Client) {
		nsa, err = c.GetServiceAuthorization(&GetServiceAuthorizationInput{
			ID: sa.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if nsa.Service.ID != TestDeliveryServiceID {
		t.Errorf("bad service id: %v", nsa.Service)
	}

	// Update
	var usa *ServiceAuthorization
	Record(t, fixtureBase+"update", func(c *Client) {
		usa, err = c.UpdateServiceAuthorization(&UpdateServiceAuthorizationInput{
			ID:         sa.ID,
			Permission: "purge_select",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if usa.Service.ID != TestDeliveryServiceID {
		t.Errorf("bad service id: %v", usa.Service)
	}
	if usa.Permission != "purge_select" {
		t.Errorf("bad permission: %v", usa.Permission)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
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
	_, err = TestClient.GetServiceAuthorization(&GetServiceAuthorizationInput{
		ID: "",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateServiceAuthorization_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateServiceAuthorization(&CreateServiceAuthorizationInput{
		Service: &SAService{ID: ""},
		User:    &SAUser{ID: ""},
	})
	if !errors.Is(err, ErrMissingServiceAuthorizationsService) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateServiceAuthorization(&CreateServiceAuthorizationInput{
		Service: &SAService{ID: "my-service-id"},
		User:    &SAUser{ID: ""},
	})
	if !errors.Is(err, ErrMissingServiceAuthorizationsUser) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateServiceAuthorization_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateServiceAuthorization(&UpdateServiceAuthorizationInput{
		ID:         "",
		Permission: "",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateServiceAuthorization(&UpdateServiceAuthorizationInput{
		ID:         "my-service-authorization-id",
		Permission: "",
	})
	if !errors.Is(err, ErrMissingPermission) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteServiceAuthorization_validation(t *testing.T) {
	err := TestClient.DeleteServiceAuthorization(&DeleteServiceAuthorizationInput{
		ID: "",
	})
	if !errors.Is(err, ErrMissingID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_listServiceAuthorizations_formatFilters(t *testing.T) {
	cases := []struct {
		remote *ListServiceAuthorizationsInput
		local  map[string]string
	}{
		{
			remote: &ListServiceAuthorizationsInput{
				PageSize:   2,
				PageNumber: 2,
			},
			local: map[string]string{
				jsonapi.QueryParamPageSize:   "2",
				jsonapi.QueryParamPageNumber: "2",
			},
		},
	}
	for _, c := range cases {
		out := c.remote.formatFilters()
		if !reflect.DeepEqual(out, c.local) {
			t.Fatalf("Error matching:\nexpected: %#v\n     got: %#v", c.local, out)
		}
	}
}
