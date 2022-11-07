package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// User represents a user of the Fastly API and web interface.
type User struct {
	CreatedAt              *time.Time `mapstructure:"created_at"`
	CustomerID             string     `mapstructure:"customer_id"`
	DeletedAt              *time.Time `mapstructure:"deleted_at"`
	EmailHash              string     `mapstructure:"email_hash"`
	ID                     string     `mapstructure:"id"`
	LimitServices          bool       `mapstructure:"limit_services"`
	Locked                 bool       `mapstructure:"locked"`
	Login                  string     `mapstructure:"login"`
	Name                   string     `mapstructure:"name"`
	RequireNewPassword     bool       `mapstructure:"require_new_password"`
	Role                   string     `mapstructure:"role"`
	TwoFactorAuthEnabled   bool       `mapstructure:"two_factor_auth_enabled"`
	TwoFactorSetupRequired bool       `mapstructure:"two_factor_setup_required"`
	UpdatedAt              *time.Time `mapstructure:"updated_at"`
}

// usersByLogin is a sortable list of users.
type usersByName []*User

// Len implement the sortable interface.
func (s usersByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s usersByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s usersByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListCustomerUsersInput is used as input to the ListCustomerUsers function.
type ListCustomerUsersInput struct {
	// CustomerID is an alphanumeric string identifying the customer (required).
	CustomerID string
}

// ListCustomerUsers retrieves all resources.
func (c *Client) ListCustomerUsers(i *ListCustomerUsersInput) ([]*User, error) {
	if i.CustomerID == "" {
		return nil, ErrMissingCustomerID
	}

	path := fmt.Sprintf("/customer/%s/users", i.CustomerID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u []*User
	if err := decodeBodyMap(resp.Body, &u); err != nil {
		return nil, err
	}
	sort.Stable(usersByName(u))
	return u, nil
}

// GetCurrentUser retrieves the user information for the authenticated user.
func (c *Client) GetCurrentUser() (*User, error) {
	resp, err := c.Get("/current_user", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u *User
	if err := decodeBodyMap(resp.Body, &u); err != nil {
		return nil, err
	}

	return u, nil
}

// GetUserInput is used as input to the GetUser function.
type GetUserInput struct {
	// ID is an alphanumeric string identifying the user (required).
	ID string
}

// GetUser retrieves the specified resource.
//
//If no user exists for the given id, the API returns a 404 response.
func (c *Client) GetUser(i *GetUserInput) (*User, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/user/%s", i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u *User
	if err := decodeBodyMap(resp.Body, &u); err != nil {
		return nil, err
	}

	return u, nil
}

// CreateUserInput is used as input to the CreateUser function.
type CreateUserInput struct {
	// Login is the login associated with the user (typically, an email address).
	Login *string `url:"login,omitempty"`
	// Name is the real life name of the user.
	Name *string `url:"name,omitempty"`
	// Role is the permissions role assigned to the user. Can be user, billing, engineer, or superuser.
	Role *string `url:"role,omitempty"`
}

// CreateUser creates a new resource.
func (c *Client) CreateUser(i *CreateUserInput) (*User, error) {
	resp, err := c.PostForm("/user", i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u *User
	if err := decodeBodyMap(resp.Body, &u); err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateUserInput is used as input to the UpdateUser function.
type UpdateUserInput struct {
	// ID is an alphanumeric string identifying the user (required).
	ID string `url:"-"`
	// Name is the real life name of the user.
	Name *string `url:"name,omitempty"`
	// Role is the permissions role assigned to the user. Can be user, billing, engineer, or superuser.
	Role *string `url:"role,omitempty"`
}

// UpdateUser updates the specified resource.
func (c *Client) UpdateUser(i *UpdateUserInput) (*User, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/user/%s", i.ID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u *User
	if err := decodeBodyMap(resp.Body, &u); err != nil {
		return nil, err
	}
	return u, nil
}

// DeleteUserInput is used as input to the DeleteUser function.
type DeleteUserInput struct {
	// ID is an alphanumeric string identifying the user (required).
	ID string
}

// DeleteUser deletes the specified resource.
func (c *Client) DeleteUser(i *DeleteUserInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/user/%s", i.ID)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}

// ResetUserPasswordInput is used as input to the ResetUserPassword function.
type ResetUserPasswordInput struct {
	// Login is the login associated with the user and is typically an email address (required).
	Login string
}

// ResetUserPassword revokes a specific token by its ID.
func (c *Client) ResetUserPassword(i *ResetUserPasswordInput) error {
	if i.Login == "" {
		return ErrMissingLogin
	}

	path := fmt.Sprintf("/user/%s/password/request_reset", url.PathEscape(i.Login))
	resp, err := c.Post(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
