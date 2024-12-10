package fastly

import (
	"strconv"
	"time"
)

// HTTP3 represents a response from the Fastly API.
type HTTP3 struct {
	CreatedAt       *time.Time `mapstructure:"created_at" json:"created_at"`
	DeletedAt       *time.Time `mapstructure:"deleted_at" json:"deleted_at"`
	FeatureRevision *int       `mapstructure:"feature_revision" json:"feature_revision"`
	ServiceID       *string    `mapstructure:"service_id" json:"service_id"`
	ServiceVersion  *int       `mapstructure:"version" json:"version"`
	UpdatedAt       *time.Time `mapstructure:"updated_at" json:"updated_at"`
}

// GetHTTP3Input is used as input to the GetHTTP3 function.
type GetHTTP3Input struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetHTTP3 retrieves the specified resource.
func (c *Client) GetHTTP3(i *GetHTTP3Input) (*HTTP3, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "http3")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *HTTP3
	if err := DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}

	return h, nil
}

// EnableHTTP3Input is used as input to the EnableHTTP3 function.
type EnableHTTP3Input struct {
	// FeatureRevision is the revision number of the HTTP/3 feature implementation.
	FeatureRevision *int `url:"feature_revision,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// EnableHTTP3 creates a new resource.
func (c *Client) EnableHTTP3(i *EnableHTTP3Input) (*HTTP3, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "http3")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var http3 *HTTP3
	if err := DecodeBodyMap(resp.Body, &http3); err != nil {
		return nil, err
	}
	return http3, nil
}

// DisableHTTP3Input is the input parameter to the DisableHTTP3 function.
type DisableHTTP3Input struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DisableHTTP3 deletes the specified resource.
func (c *Client) DisableHTTP3(i *DisableHTTP3Input) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "http3")
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
