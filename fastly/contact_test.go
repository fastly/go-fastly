package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_Contacts(t *testing.T) {
	t.Parallel()

	fixtureBase := "contacts/"

	// Need a customer ID; fetch it from the current user.
	var (
		err        error
		customerID string
	)
	Record(t, fixtureBase+"get_current_user", func(c *Client) {
		var u *User
		u, err = c.GetCurrentUser(context.TODO())
		if err == nil && u != nil && u.CustomerID != nil {
			customerID = *u.CustomerID
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if customerID == "" {
		t.Fatal("missing customer id from current user")
	}

	// Fastly refuses to delete the last contact of a given type, so we
	// create a guard contact first that we leave in place for the duration
	// of the test (and best-effort delete during cleanup).
	//
	// NOTE: When recreating the fixtures, update both emails.
	guardEmail := "go-fastly-test+contact-guard+20260522@example.com"
	email := "go-fastly-test+contact+20260522@example.com"

	var guard *Contact
	Record(t, fixtureBase+"create_guard", func(c *Client) {
		guard, err = c.CreateContact(context.TODO(), &CreateContactInput{
			CustomerID:  customerID,
			ContactType: "emergency",
			Name:        "guard contact",
			Email:       guardEmail,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create the contact under test.
	var co *Contact
	Record(t, fixtureBase+"create", func(c *Client) {
		co, err = c.CreateContact(context.TODO(), &CreateContactInput{
			CustomerID:  customerID,
			ContactType: "emergency",
			Name:        "test contact",
			Email:       email,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted (best-effort; guard deletion may fail if it is the
	// last emergency contact on the account, which is fine).
	defer func() {
		Record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteContact(context.TODO(), &DeleteContactInput{
				CustomerID: customerID,
				ContactID:  co.ContactID,
			})
			_ = c.DeleteContact(context.TODO(), &DeleteContactInput{
				CustomerID: customerID,
				ContactID:  guard.ContactID,
			})
		})
	}()

	if co.ContactID == "" {
		t.Errorf("bad contact id: %+v", co)
	}
	if co.Email != email {
		t.Errorf("bad email: %q", co.Email)
	}

	// List
	var cs []*Contact
	Record(t, fixtureBase+"list", func(c *Client) {
		cs, err = c.ListContacts(context.TODO(), &ListContactsInput{
			CustomerID: customerID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cs) < 1 {
		t.Errorf("bad contacts: %v", cs)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteContact(context.TODO(), &DeleteContactInput{
			CustomerID: customerID,
			ContactID:  co.ContactID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListContacts_validation(t *testing.T) {
	_, err := TestClient.ListContacts(context.TODO(), &ListContactsInput{
		CustomerID: "",
	})
	if !errors.Is(err, ErrMissingCustomerID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateContact_validation(t *testing.T) {
	_, err := TestClient.CreateContact(context.TODO(), &CreateContactInput{
		CustomerID: "",
	})
	if !errors.Is(err, ErrMissingCustomerID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteContact_validation(t *testing.T) {
	err := TestClient.DeleteContact(context.TODO(), &DeleteContactInput{
		CustomerID: "",
		ContactID:  "abc",
	})
	if !errors.Is(err, ErrMissingCustomerID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteContact(context.TODO(), &DeleteContactInput{
		CustomerID: "abc",
		ContactID:  "",
	})
	if !errors.Is(err, ErrMissingContactID) {
		t.Errorf("bad error: %s", err)
	}
}
