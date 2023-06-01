package fastly

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/jsonapi"
)

// TLSMutualAuthentication represents a mutual authentication certificate.
type TLSMutualAuthentication struct {
	Activations []*TLSActivation `jsonapi:"relation,tls_activations"`
	CreatedAt   *time.Time       `jsonapi:"attr,created_at,iso8601"`
	Enforced    bool             `jsonapi:"attr,enforced"`
	ID          string           `jsonapi:"primary,mutual_authentication"`
	Name        string           `jsonapi:"attr,name"`
	UpdatedAt   *time.Time       `jsonapi:"attr,updated_at,iso8601"`
}

// ListTLSMutualAuthenticationsInput is used as input to the Client.ListTLSMutualAuthentication function.
type ListTLSMutualAuthenticationsInput struct {
	// Include is a list of related objects to include (optional).
	Include []string
	// PageNumber is the required page index for pagination.
	PageNumber int
	// PageSize is the number of records per page.
	PageSize int
}

// formatFilters converts user input into query parameters for filtering.
func (i *ListTLSMutualAuthenticationsInput) formatFilters() map[string]string {
	result := map[string]string{}
	if len(i.Include) > 0 {
		result["include"] = strings.Join(i.Include, ",")
	}
	if i.PageSize > 0 {
		result["page[size]"] = strconv.Itoa(i.PageSize)
	}
	if i.PageNumber > 0 {
		result["page[number]"] = strconv.Itoa(i.PageNumber)
	}
	return result
}

// ListTLSMutualAuthentication retrieves all resources.
func (c *Client) ListTLSMutualAuthentication(i *ListTLSMutualAuthenticationsInput) ([]*TLSMutualAuthentication, error) {
	p := "/tls/mutual_authentications"
	filters := &RequestOptions{
		Params: i.formatFilters(),
		Headers: map[string]string{
			"Accept": "application/vnd.api+json", // this is required otherwise the filters don't work
		},
	}

	r, err := c.Get(p, filters)
	if err != nil {
		return nil, err
	}

	data, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(TLSMutualAuthentication)))
	if err != nil {
		return nil, err
	}

	so := make([]*TLSMutualAuthentication, len(data))
	for i := range data {
		typed, ok := data[i].(*TLSMutualAuthentication)
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T", data[i])
		}
		so[i] = typed
	}

	return so, nil
}

// GetTLSMutualAuthenticationInput is used as input to the GetTLSMutualAuthentication function.
type GetTLSMutualAuthenticationInput struct {
	// ID is an alphanumeric string identifying a mutual authentication (required).
	ID string
	// Include is a comma-separated list of related objects to include (optional).
	Include string
}

// GetTLSMutualAuthentication retrieves the specified resource.
func (c *Client) GetTLSMutualAuthentication(i *GetTLSMutualAuthenticationInput) (*TLSMutualAuthentication, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	var ro *RequestOptions
	if i.Include != "" {
		ro = &RequestOptions{
			Params: map[string]string{
				"include": i.Include,
			},
			Headers: map[string]string{
				"Accept": "application/vnd.api+json", // this is required otherwise the filters don't work
			},
		}
	}

	p := fmt.Sprintf("/tls/mutual_authentications/%s", i.ID)

	r, err := c.Get(p, ro)
	if err != nil {
		return nil, err
	}

	var o TLSMutualAuthentication
	if err := jsonapi.UnmarshalPayload(r.Body, &o); err != nil {
		return nil, err
	}

	return &o, nil
}

// CreateTLSMutualAuthenticationInput is used as input to the CreateTLSMutualAuthentication function.
type CreateTLSMutualAuthenticationInput struct {
	// CertBundle is one or more certificates. Enter each individual certificate blob on a new line. Must be PEM-formatted (required).
	CertBundle string `jsonapi:"attr,cert_bundle"`
	// Enforced determines whether Mutual TLS will fail closed (enforced) or fail open.
	Enforced bool `jsonapi:"attr,enforced"`
	// ID should not be set (it's internally used to help marshal the JSONAPI request data).
	ID string `jsonapi:"primary,mutual_authentication"`
	// Name is a custom name for your mutual authentication (optional).
	Name string `jsonapi:"attr,name,omitempty"`
}

// CreateTLSMutualAuthentication creates a new resource.
func (c *Client) CreateTLSMutualAuthentication(i *CreateTLSMutualAuthenticationInput) (*TLSMutualAuthentication, error) {
	if i.CertBundle == "" {
		return nil, ErrMissingCertBundle
	}

	p := "/tls/mutual_authentications"

	r, err := c.PostJSONAPI(p, i, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var o TLSMutualAuthentication
	if err := jsonapi.UnmarshalPayload(r.Body, &o); err != nil {
		return nil, err
	}

	return &o, nil
}

// UpdateTLSMutualAuthenticationInput is used as input to the UpdateTLSMutualAuthentication function.
type UpdateTLSMutualAuthenticationInput struct {
	// CertBundle is one or more certificates. Enter each individual certificate blob on a new line. Must be PEM-formatted (required).
	CertBundle string `jsonapi:"attr,cert_bundle"`
	// Enforced determines whether Mutual TLS will fail closed (enforced) or fail open.
	Enforced bool `jsonapi:"attr,enforced"`
	// ID is an alphanumeric string identifying a mutual authentication (required).
	ID string `jsonapi:"primary,mutual_authentication"`
	// Name is a custom name for your mutual authentication (optional).
	Name string `jsonapi:"attr,name,omitempty"`
}

// UpdateTLSMutualAuthentication updates the specified resource.
func (c *Client) UpdateTLSMutualAuthentication(i *UpdateTLSMutualAuthenticationInput) (*TLSMutualAuthentication, error) {
	if i.CertBundle == "" {
		return nil, ErrMissingCertBundle
	}
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/tls/mutual_authentications/%s", i.ID)
	resp, err := c.PatchJSONAPI(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var o TLSMutualAuthentication
	if err := jsonapi.UnmarshalPayload(resp.Body, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// DeleteTLSMutualAuthenticationInput used for deleting a certificate.
type DeleteTLSMutualAuthenticationInput struct {
	// ID is an alphanumeric string identifying a mutual authentication (required).
	ID string
}

// DeleteTLSMutualAuthentication deletes the specified resource.
func (c *Client) DeleteTLSMutualAuthentication(i *DeleteTLSMutualAuthenticationInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/tls/mutual_authentications/%s", i.ID)
	_, err := c.Delete(path, nil)
	return err
}
