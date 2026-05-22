package fastly

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/jsonapi"
)

// Contact represents a customer contact.
type Contact struct {
	ContactID   string     `jsonapi:"primary,customer_contact"`
	UserID      string     `jsonapi:"attr,user_id,omitempty"`
	ContactType string     `jsonapi:"attr,contact_type,omitempty"`
	Name        string     `jsonapi:"attr,name,omitempty"`
	FirstName   string     `jsonapi:"attr,first_name,omitempty"`
	LastName    string     `jsonapi:"attr,last_name,omitempty"`
	Email       string     `jsonapi:"attr,email,omitempty"`
	Phone       string     `jsonapi:"attr,phone,omitempty"`
	CreatedAt   *time.Time `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt   *time.Time `jsonapi:"attr,updated_at,iso8601"`
	DeletedAt   *time.Time `jsonapi:"attr,deleted_at,iso8601"`
}

// ListContactsInput is used as input to the ListContacts function.
type ListContactsInput struct {
	// CustomerID is an alphanumeric string identifying the customer (required).
	CustomerID string
}

// ListContacts retrieves all contacts for the given customer.
func (c *Client) ListContacts(ctx context.Context, i *ListContactsInput) ([]*Contact, error) {
	if i.CustomerID == "" {
		return nil, ErrMissingCustomerID
	}

	path := ToSafeURL("customer", i.CustomerID, "contacts")

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = jsonapi.MediaType

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Contact)))
	if err != nil {
		return nil, err
	}

	cs := make([]*Contact, len(data))
	for idx := range data {
		typed, ok := data[idx].(*Contact)
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T", data[idx])
		}
		cs[idx] = typed
	}

	return cs, nil
}

// CreateContactInput is used as input to the CreateContact function.
type CreateContactInput struct {
	// CustomerID is an alphanumeric string identifying the customer (required).
	CustomerID string
	// UserID is an alphanumeric string identifying a user. Required when not
	// providing Email and Name.
	UserID string
	// ContactType is the type of contact. One of: primary, billing, technical,
	// security, emergency.
	ContactType string
	// Name is the name of this contact, when not referencing an existing user.
	Name string
	// FirstName is the first name of this contact, when not referencing an
	// existing user.
	FirstName string
	// LastName is the last name of this contact, when not referencing an
	// existing user.
	LastName string
	// Email is the email of this contact, when not referencing an existing user.
	Email string
	// Phone is the contact's phone number. Required for the primary, technical,
	// and security contact types.
	Phone string
}

// createContactPayload is the JSON:API marshaling shape for CreateContact.
type createContactPayload struct {
	ContactID   string `jsonapi:"primary,customer_contact"`
	UserID      string `jsonapi:"attr,user_id,omitempty"`
	ContactType string `jsonapi:"attr,contact_type,omitempty"`
	Name        string `jsonapi:"attr,name,omitempty"`
	FirstName   string `jsonapi:"attr,first_name,omitempty"`
	LastName    string `jsonapi:"attr,last_name,omitempty"`
	Email       string `jsonapi:"attr,email,omitempty"`
	Phone       string `jsonapi:"attr,phone,omitempty"`
}

// CreateContact creates a new contact.
func (c *Client) CreateContact(ctx context.Context, i *CreateContactInput) (*Contact, error) {
	if i.CustomerID == "" {
		return nil, ErrMissingCustomerID
	}

	path := ToSafeURL("customer", i.CustomerID, "contacts")

	payload := &createContactPayload{
		UserID:      i.UserID,
		ContactType: i.ContactType,
		Name:        i.Name,
		FirstName:   i.FirstName,
		LastName:    i.LastName,
		Email:       i.Email,
		Phone:       i.Phone,
	}

	resp, err := c.PostJSONAPI(ctx, path, payload, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var co Contact
	if err := jsonapi.UnmarshalPayload(resp.Body, &co); err != nil {
		return nil, err
	}
	return &co, nil
}

// DeleteContactInput is used as input to the DeleteContact function.
type DeleteContactInput struct {
	// CustomerID is an alphanumeric string identifying the customer (required).
	CustomerID string
	// ContactID is an alphanumeric string identifying the contact (required).
	ContactID string
}

// DeleteContact deletes the specified contact.
func (c *Client) DeleteContact(ctx context.Context, i *DeleteContactInput) error {
	if i.CustomerID == "" {
		return ErrMissingCustomerID
	}
	if i.ContactID == "" {
		return ErrMissingContactID
	}

	path := ToSafeURL("customer", i.CustomerID, "contacts", i.ContactID)

	resp, err := c.Delete(ctx, path, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
