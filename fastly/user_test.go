package fastly

import (
	"errors"
	"testing"
)

func TestClient_UsersCurrent(t *testing.T) {
	t.Parallel()

	var err error
	var u *User
	Record(t, "users/get_current_user", func(c *Client) {
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
	login := "go-fastly-test+user+20221104@example.com"

	// Create
	//
	// NOTE: When recreating the fixtures, update the login.
	var err error
	var u *User
	Record(t, fixtureBase+"create", func(c *Client) {
		u, err = c.CreateUser(&CreateUserInput{
			Login: ToPointer(login),
			Name:  ToPointer("test user"),
			Role:  ToPointer("engineer"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteUser(&DeleteUserInput{
				UserID: *u.UserID,
			})
		})
	}()

	if *u.Login != login {
		t.Errorf("bad login: %v", *u.Login)
	}

	if *u.Name != "test user" {
		t.Errorf("bad name: %v", *u.Name)
	}

	if *u.Role != "engineer" {
		t.Errorf("bad role: %v", *u.Role)
	}

	// List
	var us []*User
	Record(t, fixtureBase+"list", func(c *Client) {
		us, err = c.ListCustomerUsers(&ListCustomerUsersInput{
			CustomerID: *u.CustomerID,
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
	Record(t, fixtureBase+"get", func(c *Client) {
		nu, err = c.GetUser(&GetUserInput{
			UserID: *u.UserID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *u.Name != *nu.Name {
		t.Errorf("bad name: %q (%q)", *u.Name, *nu.Name)
	}

	// Update
	var uu *User
	Record(t, fixtureBase+"update", func(c *Client) {
		uu, err = c.UpdateUser(&UpdateUserInput{
			UserID: *u.UserID,
			Name:   ToPointer("updated user"),
			Role:   ToPointer("superuser"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uu.Name != "updated user" {
		t.Errorf("bad name: %q", *uu.Name)
	}
	if *uu.Role != "superuser" {
		t.Errorf("bad role: %q", *uu.Role)
	}

	// Reset Password
	//
	// NOTE: This integration test can fail due to reCAPTCHA.
	// Which means you might have to manually correct the fixtures 😬
	Record(t, fixtureBase+"reset_password", func(c *Client) {
		err = c.ResetUserPassword(&ResetUserPasswordInput{
			Login: *uu.Login,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteUser(&DeleteUserInput{
			UserID: *u.UserID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListCustomerUsers_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListCustomerUsers(&ListCustomerUsersInput{
		CustomerID: "",
	})
	if !errors.Is(err, ErrMissingCustomerID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetUser_validation(t *testing.T) {
	var err error
	_, err = TestClient.GetUser(&GetUserInput{
		UserID: "",
	})
	if !errors.Is(err, ErrMissingUserID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateUser_validation(t *testing.T) {
	var err error
	_, err = TestClient.UpdateUser(&UpdateUserInput{
		UserID: "",
	})
	if !errors.Is(err, ErrMissingUserID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteUser_validation(t *testing.T) {
	err := TestClient.DeleteUser(&DeleteUserInput{
		UserID: "",
	})
	if !errors.Is(err, ErrMissingUserID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ResetUser_validation(t *testing.T) {
	err := TestClient.ResetUserPassword(&ResetUserPasswordInput{
		Login: "",
	})
	if !errors.Is(err, ErrMissingLogin) {
		t.Errorf("bad error: %s", err)
	}
}
