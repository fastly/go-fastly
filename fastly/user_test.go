package fastly

import "testing"

func TestClient_ListUsers(t *testing.T) {
	t.Parallel()

	var users []*User
	var err error
	record(t, "users/list", func(c *Client) {
		users, err = c.ListUsers()
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) < 1 {
		t.Errorf("bad users: %v", users)
	}
}

func TestClient_ListCustomerUsers(t *testing.T) {
	t.Parallel()

	var users []*User
	var err error
	record(t, "users/list_customer", func(c *Client) {
		users, err = c.ListCustomerUsers(&ListCustomerUsersInput{
			ID: "CUsT0m3rID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) < 1 {
		t.Errorf("bad users: %v", users)
	}
}

func TestClient_GetCurrentUser(t *testing.T) {
	t.Parallel()

	var user *User
	var err error
	record(t, "users/get_current_user", func(c *Client) {
		user, err = c.GetCurrentUser()
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", user)
}

func TestClient_GetUser(t *testing.T) {
	t.Parallel()

	input := &GetUserInput{
		ID: "T3hr33eE3thr3EEtHr3e3",
	}

	var user *User
	var err error
	record(t, "users/get", func(c *Client) {
		user, err = c.GetUser(input)
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", user)
}

func TestClient_CreateUser(t *testing.T) {
	t.Parallel()

	input := &CreateUserInput{
		Login: "new+engineer@example.com",
		Name:  "new engineer",
		Role:  "engineer",
	}

	var user *User
	var err error
	record(t, "users/create", func(c *Client) {
		user, err = c.CreateUser(input)
	})
	if err != nil {
		t.Fatal(err)
	}

	if user.Login != input.Login {
		t.Errorf("returned invalid login, got %s, want %s", user.Login, input.Login)
	}

	if user.Name != input.Name {
		t.Errorf("returned invalid name, got %s, want %s", user.Name, input.Name)
	}

	if user.Role != input.Role {
		t.Errorf("returned invalid role, got %s, want %s", user.Role, input.Role)
	}
}

func TestClient_UpdateUser(t *testing.T) {
	t.Parallel()

	input := &UpdateUserInput{
		ID:   "T3hr33eE3thr3EEtHr3e3",
		Name: "Superuser Three",
		Role: "superuser",
	}

	var user *User
	var err error
	record(t, "users/update", func(c *Client) {
		user, err = c.UpdateUser(input)
	})
	if err != nil {
		t.Fatal(err)
	}

	if user.Name != input.Name {
		t.Errorf("returned invalid name, got %s, want %s", user.Name, input.Name)
	}

	if user.Role != input.Role {
		t.Errorf("returned invalid role, got %s, want %s", user.Role, input.Role)
	}
}

func TestClient_DeleteUser(t *testing.T) {
	t.Parallel()

	input := &DeleteUserInput{
		ID: "6SI6xsIX66S66iX6ixSIx",
	}

	var err error
	record(t, "users/delete", func(c *Client) {
		err = c.DeleteUser(input)
	})
	if err != nil {
		t.Fatal(err)
	}
}
