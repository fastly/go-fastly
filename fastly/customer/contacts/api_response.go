package contacts

import "time"

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
