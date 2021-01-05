package fastly

import (
	"testing"
)

func TestClient_GetCurrentUser(t *testing.T) {
	t.Parallel()

	var err error
	var u *User
	record(t, "users/get_current_user", func(c *Client) {
		u, err = c.GetCurrentUser()
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", u)
}

func TestClient_Users(t *testing.T) {
	t.Parallel()

	fixtureBase := "users/"

	// Create
	var err error
	var u *User
	record(t, fixtureBase+"create", func(c *Client) {
		u, err = c.CreateUser(&CreateUserInput{
			Login: "go-fastly-test+user1@example.com",
			Name:  "test user",
			Role:  "engineer",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeleteUser(&DeleteUserInput{
				ID: u.ID,
			})
		})
	}()

	if u.Login != "go-fastly-test+user1@example.com" {
		t.Errorf("bad login: %v", u.Login)
	}

	if u.Name != "test user" {
		t.Errorf("bad name: %v", u.Name)
	}

	if u.Role != "engineer" {
		t.Errorf("bad role: %v", u.Role)
	}

	// List
	var us []*User
	record(t, fixtureBase+"list", func(c *Client) {
		us, err = c.ListCustomerUsers(&ListCustomerUsersInput{
			CustomerID: u.CustomerID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) < 1 {
		t.Errorf("bad users: %v", us)
	}

	// Get
	var nu *User
	record(t, fixtureBase+"get", func(c *Client) {
		nu, err = c.GetUser(&GetUserInput{
			ID: u.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != nu.Name {
		t.Errorf("bad name: %q (%q)", u.Name, nu.Name)
	}

	// Update
	var uu *User
	record(t, fixtureBase+"update", func(c *Client) {
		uu, err = c.UpdateUser(&UpdateUserInput{
			ID:   u.ID,
			Name: String("updated user"),
			Role: String("superuser"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uu.Name != "updated user" {
		t.Errorf("bad name: %q", uu.Name)
	}
	if uu.Role != "superuser" {
		t.Errorf("bad role: %q", uu.Role)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteUser(&DeleteUserInput{
			ID: u.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateUser_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateUser(&CreateUserInput{
		Login: "",
	})
	if err != ErrMissingLogin {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateUser(&CreateUserInput{
		Login: "new+user@example.com",
		Name:  "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ListCustomerUsers_validation(t *testing.T) {
	var err error
	_, err = testClient.ListCustomerUsers(&ListCustomerUsersInput{
		CustomerID: "",
	})
	if err != ErrMissingCustomerID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetUser_validation(t *testing.T) {
	var err error
	_, err = testClient.GetUser(&GetUserInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateUser_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateUser(&UpdateUserInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteUser_validation(t *testing.T) {
	err := testClient.DeleteUser(&DeleteUserInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
