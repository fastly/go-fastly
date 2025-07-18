package fastly

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/jsonapi"
)

// TLSSubscription represents a managed TLS certificate.
type TLSSubscription struct {
	Authorizations       []*TLSAuthorizations          `jsonapi:"relation,tls_authorizations"`
	CertificateAuthority string                        `jsonapi:"attr,certificate_authority"`
	Certificates         []*TLSSubscriptionCertificate `jsonapi:"relation,tls_certificates"`
	CommonName           *TLSDomain                    `jsonapi:"relation,common_name"`
	Configuration        *TLSConfiguration             `jsonapi:"relation,tls_configuration"`
	CreatedAt            *time.Time                    `jsonapi:"attr,created_at,iso8601"`
	Domains              []*TLSDomain                  `jsonapi:"relation,tls_domains"`
	ID                   string                        `jsonapi:"primary,tls_subscription"`
	State                string                        `jsonapi:"attr,state"`
	UpdatedAt            *time.Time                    `jsonapi:"attr,updated_at,iso8601"`
}

// TLSSubscriptionCertificate represents a subscription certificate.
type TLSSubscriptionCertificate struct {
	ID                 string     `jsonapi:"primary,tls_certificate"`
	CreatedAt          *time.Time `jsonapi:"attr,created_at,iso8601"`
	IssuedTo           string     `jsonapi:"attr,issued_to"`
	Issuer             string     `jsonapi:"attr,issuer"`
	Name               string     `jsonapi:"attr,name"`
	NotAfter           *time.Time `jsonapi:"attr,not_after,iso8601"`
	NotBefore          *time.Time `jsonapi:"attr,not_before,iso8601"`
	Replace            bool       `jsonapi:"attr,replace"`
	SerialNumber       string     `jsonapi:"attr,serial_number"`
	SignatureAlgorithm string     `jsonapi:"attr,signature_algorithm"`
	UpdatedAt          *time.Time `jsonapi:"attr,updated_at,iso8601"`
}

// TLSAuthorizations gives information needed to verify domain ownership in
// order to enable a TLSSubscription.
type TLSAuthorizations struct {
	// Challenges ...
	// See https://github.com/google/jsonapi/pull/99
	// WARNING: Nested structs only work with values, not pointers.
	Challenges []TLSChallenge            `jsonapi:"attr,challenges"`
	CreatedAt  *time.Time                `jsonapi:"attr,created_at,iso8601,omitempty"`
	ID         string                    `jsonapi:"primary,tls_authorization"`
	State      string                    `jsonapi:"attr,state,omitempty"`
	UpdatedAt  *time.Time                `jsonapi:"attr,updated_at,iso8601,omitempty"`
	Warnings   []TLSAuthorizationWarning `jsonapi:"attr,warnings,omitempty"`
}

// TLSAuthorizationWarning indicates possible issues with the TLS configuration.
type TLSAuthorizationWarning struct {
	Type         string `jsonapi:"attr,type"`
	Instructions string `jsonapi:"attr,instructions"`
}

// TLSChallenge represents a DNS record to be added for a specific type of
// domain ownership challenge.
type TLSChallenge struct {
	RecordName string   `jsonapi:"attr,record_name"`
	RecordType string   `jsonapi:"attr,record_type"`
	Type       string   `jsonapi:"attr,type"`
	Values     []string `jsonapi:"attr,values"`
}

// ListTLSSubscriptionsInput is used as input to the ListTLSSubscriptions
// function.
type ListTLSSubscriptionsInput struct {
	// Limit the returned subscriptions to those that have currently active orders. Permitted values: true.
	FilterActiveOrders bool
	// Limit the returned subscriptions by state. Valid values are pending, processing, issued, and renewing. Accepts parameters: not (e.g., filter[state][not]=renewing).
	FilterState string
	// Limit the returned subscriptions to those that include the specific domain.
	FilterTLSDomainsID string
	// Include related objects. Optional, comma-separated values. Permitted values: tls_authorizations, tls_authorizations.globalsign_email_challenge, tls_authorizations.self_managed_http_challenge, and tls_certificates.
	Include string
	// Current page.
	PageNumber int
	// Number of records per page.
	PageSize int
	// The order in which to list the results by creation date. Accepts created_at (ascending sort order) or -created_at (descending).
	Sort string
}

// formatFilters converts user input into query parameters for filtering.
func (s *ListTLSSubscriptionsInput) formatFilters() map[string]string {
	result := map[string]string{}
	pairings := map[string]any{
		"filter[has_active_order]":   s.FilterActiveOrders,
		"filter[state]":              s.FilterState,
		"filter[tls_domains.id]":     s.FilterTLSDomainsID,
		"include":                    s.Include,
		jsonapi.QueryParamPageNumber: s.PageNumber,
		jsonapi.QueryParamPageSize:   s.PageSize,
		"sort":                       s.Sort,
	}

	for key, v := range pairings {
		switch value := v.(type) {
		case bool:
			// NOTE: The API currently has a bug where the presence of the
			// has_active_order filter will cause the response to include
			// subscriptions with an active order, even if the filter value itself was
			// set to `false`. This is considered a bug and the Fastly API team are
			// aware of the issue. For now, go-fastly will omit setting the filter
			// unless the key includes has_active_order and the value is explicitly
			// set to `true`.
			if (key == "filter[has_active_order]" && value) || key != "filter[has_active_order]" {
				result[key] = strconv.FormatBool(value)
			}
		case string:
			if value != "" {
				result[key] = value
			}
		case int:
			if value != 0 {
				result[key] = strconv.Itoa(value)
			}
		}
	}
	return result
}

// ListTLSSubscriptions retrieves all resources.
func (c *Client) ListTLSSubscriptions(ctx context.Context, i *ListTLSSubscriptionsInput) ([]*TLSSubscription, error) {
	requestOptions := CreateRequestOptions()
	requestOptions.Params = i.formatFilters()
	requestOptions.Headers["Accept"] = jsonapi.MediaType // this is required otherwise the filters don't work

	resp, err := c.Get(ctx, "/tls/subscriptions", requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(TLSSubscription)))
	if err != nil {
		return nil, err
	}

	// Convert slice of any to a slice of TLSSubscription structs
	subscriptions := make([]*TLSSubscription, len(data))
	for i := range data {
		typed, ok := data[i].(*TLSSubscription)
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T", data[i])
		}
		subscriptions[i] = typed
	}

	return subscriptions, nil
}

// CreateTLSSubscriptionInput is used as input to the CreateTLSSubscription
// function.
type CreateTLSSubscriptionInput struct {
	// CertificateAuthority is the entity that issues and certifies the TLS certificates for your subscription. Valid values are lets-encrypt or globalsign.
	CertificateAuthority string `jsonapi:"attr,certificate_authority,omitempty"`
	// CommonName is the common name associated with the subscription generated by Fastly TLS. Must be included in Domains. Only the ID fields of each one need to be set.
	CommonName *TLSDomain `jsonapi:"relation,common_name,omitempty"`
	// Configuration options that apply to the enabled domains on this subscription. Only ID needs to be populated
	Configuration *TLSConfiguration `jsonapi:"relation,tls_configuration,omitempty"`
	// Domains list to enable TLS for. Only the ID fields of each one need to be set.
	Domains []*TLSDomain `jsonapi:"relation,tls_domain"`
	// ID value is ignored and should not be set, needed to make JSONAPI work correctly.
	ID string `jsonapi:"primary,tls_subscription"`
}

// CreateTLSSubscription creates a new resource.
func (c *Client) CreateTLSSubscription(ctx context.Context, i *CreateTLSSubscriptionInput) (*TLSSubscription, error) {
	if len(i.Domains) == 0 {
		return nil, ErrMissingTLSDomain
	}
	if i.CommonName != nil && !domainInSlice(i.Domains, i.CommonName) {
		return nil, ErrCommonNameNotInDomains
	}

	resp, err := c.PostJSONAPI(ctx, "/tls/subscriptions", i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subscription TLSSubscription
	err = jsonapi.UnmarshalPayload(resp.Body, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// domainInSlice takes a slice of TLSDomain structs, and another TLSDomain struct to search for, returning true if any
// of the ID fields in the slice match.
func domainInSlice(haystack []*TLSDomain, needle *TLSDomain) bool {
	for _, s := range haystack {
		if s.ID == needle.ID {
			return true
		}
	}

	return false
}

// GetTLSSubscriptionInput is used as input to the GetTLSSubscription function.
type GetTLSSubscriptionInput struct {
	// ID of the TLS subscription to fetch.
	ID string
	// Include related objects. Optional, comma-separated values. Permitted values: tls_authorizations, tls_authorizations.globalsign_email_challenge, tls_authorizations.self_managed_http_challenge, and tls_certificates.
	Include *string
}

// GetTLSSubscription retrieves the specified resource.
func (c *Client) GetTLSSubscription(ctx context.Context, i *GetTLSSubscriptionInput) (*TLSSubscription, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("tls", "subscriptions", i.ID)

	requestOptions := CreateRequestOptions()
	requestOptions.Headers["Accept"] = jsonapi.MediaType // this is required otherwise the filters don't work

	if i.Include != nil {
		requestOptions.Params["include"] = *i.Include
	}

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subscription TLSSubscription
	err = jsonapi.UnmarshalPayload(resp.Body, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// UpdateTLSSubscriptionInput is used as input to the UpdateTLSSubscription
// function (Limited Availability).
type UpdateTLSSubscriptionInput struct {
	// CommonName is the common name associated with the subscription generated by Fastly TLS. Must be included in Domains. Only the ID fields of each one need to be set.
	CommonName *TLSDomain `jsonapi:"relation,common_name,omitempty"`
	// Configuration options that apply to the enabled domains on this subscription. Only ID needs to be populated
	Configuration *TLSConfiguration `jsonapi:"relation,tls_configuration,omitempty"`
	// Domains list to enable TLS for. Only the ID fields of each one need to be set.
	Domains []*TLSDomain `jsonapi:"relation,tls_domain,omitempty"`
	// Force the update to be applied, even if domains are active. Warning: can disable production traffic.
	Force bool
	// ID of the subscription to update.
	ID string `jsonapi:"primary,tls_subscription"`
}

// UpdateTLSSubscription updates the specified resource.
//
// TLS Subscriptions can only be updated in an "issued" state, and when Force=true.
func (c *Client) UpdateTLSSubscription(ctx context.Context, i *UpdateTLSSubscriptionInput) (*TLSSubscription, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := ToSafeURL("tls", "subscriptions", i.ID)

	requestOptions := CreateRequestOptions()
	if i.Force {
		requestOptions.Params["force"] = "true"
	}

	resp, err := c.PatchJSONAPI(ctx, path, i, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var subscription TLSSubscription
	err = jsonapi.UnmarshalPayload(resp.Body, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// DeleteTLSSubscriptionInput is used as input to the DeleteTLSSubscription
// function.
type DeleteTLSSubscriptionInput struct {
	// Force the subscription to be deleted, even if domains are active. Warning: can disable production traffic.
	Force bool
	// ID of the TLS subscription to delete.
	ID string
}

// DeleteTLSSubscription deletes the specified resource.
func (c *Client) DeleteTLSSubscription(ctx context.Context, i *DeleteTLSSubscriptionInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	requestOptions := CreateRequestOptions()
	if i.Force {
		requestOptions.Params["force"] = "true"
	}

	path := ToSafeURL("tls", "subscriptions", i.ID)

	ignored, err := c.Delete(ctx, path, requestOptions)
	if err != nil {
		return err
	}
	defer ignored.Body.Close()
	return nil
}
