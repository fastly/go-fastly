package contacts

import (
	"context"

	"github.com/google/jsonapi"

	"github.com/fastly/go-fastly/v15/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// CustomerID is the alphanumeric identifier of the customer (required).
	CustomerID *string
	// UserID is the alphanumeric identifier of an existing user. Required when
	// not providing Email and Name.
	UserID *string
	// ContactType is the type of contact. One of: primary, billing, technical,
	// security, emergency.
	ContactType *string
	// Name is the name of this contact, when not referencing an existing user.
	Name *string
	// FirstName is the first name of this contact, when not referencing an
	// existing user.
	FirstName *string
	// LastName is the last name of this contact, when not referencing an
	// existing user.
	LastName *string
	// Email is the email of this contact, when not referencing an existing user.
	Email *string
	// Phone is the contact's phone number. Required for the primary, technical,
	// and security contact types.
	Phone *string
}

// createPayload is the JSON:API marshaling shape for Create.
//
// It mirrors Contact but omits the fields that must not be sent in the body
// (CustomerID lives in the URL).
type createPayload struct {
	ContactID   string `jsonapi:"primary,customer_contact"`
	UserID      string `jsonapi:"attr,user_id,omitempty"`
	ContactType string `jsonapi:"attr,contact_type,omitempty"`
	Name        string `jsonapi:"attr,name,omitempty"`
	FirstName   string `jsonapi:"attr,first_name,omitempty"`
	LastName    string `jsonapi:"attr,last_name,omitempty"`
	Email       string `jsonapi:"attr,email,omitempty"`
	Phone       string `jsonapi:"attr,phone,omitempty"`
}

// Create creates a new contact for the given customer.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Contact, error) {
	if i.CustomerID == nil {
		return nil, fastly.ErrMissingCustomerID
	}

	path := fastly.ToSafeURL("customer", *i.CustomerID, "contacts")

	payload := &createPayload{
		UserID:      fastly.ToValue(i.UserID),
		ContactType: fastly.ToValue(i.ContactType),
		Name:        fastly.ToValue(i.Name),
		FirstName:   fastly.ToValue(i.FirstName),
		LastName:    fastly.ToValue(i.LastName),
		Email:       fastly.ToValue(i.Email),
		Phone:       fastly.ToValue(i.Phone),
	}

	resp, err := c.PostJSONAPI(ctx, path, payload, fastly.CreateRequestOptions())
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
