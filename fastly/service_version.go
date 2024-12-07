package fastly

import (
	"strconv"
	"time"
)

// Version represents a distinct configuration version.
type Version struct {
	Active       *bool          `mapstructure:"active"`
	Comment      *string        `mapstructure:"comment"`
	CreatedAt    *time.Time     `mapstructure:"created_at"`
	DeletedAt    *time.Time     `mapstructure:"deleted_at"`
	Deployed     *bool          `mapstructure:"deployed"`
	Locked       *bool          `mapstructure:"locked"`
	Number       *int           `mapstructure:"number"`
	ServiceID    *string        `mapstructure:"service_id"`
	Staging      *bool          `mapstructure:"staging"`
	Testing      *bool          `mapstructure:"testing"`
	UpdatedAt    *time.Time     `mapstructure:"updated_at"`
	Environments []*Environment `mapstructure:"environments"`
}

// Environment represents a distinct deployment environment.
type Environment struct {
	ServiceVersion *int64  `mapstructure:"active_version"`
	Name           *string `mapstructure:"name"`
	ServiceID      *string `mapstructure:"service_id"`
}

// ListVersionsInput is the input to the ListVersions function.
type ListVersionsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// ListVersions retrieves all resources.
func (c *Client) ListVersions(i *ListVersionsInput) ([]*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "version")

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e []*Version
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}

	return e, nil
}

// LatestVersionInput is the input to the LatestVersion function.
type LatestVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// LatestVersion retrieves the specified resource.
//
// If there are no versions, this function will return nil (but not an error).
func (c *Client) LatestVersion(i *LatestVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	slice, err := c.ListVersions(&ListVersionsInput{ServiceID: i.ServiceID})
	if err != nil {
		return nil, err
	}
	if len(slice) < 1 {
		return nil, nil
	}

	e := slice[len(slice)-1]
	return e, nil
}

// CreateVersionInput is the input to the CreateVersion function.
type CreateVersionInput struct {
	// Comment is a personal freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
}

// CreateVersion creates a new resource.
//
// This is preferred in almost all scenarios, since `Create()` creates a _blank_
// configuration where `Clone()` builds off of an existing configuration.
func (c *Client) CreateVersion(i *CreateVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "version")

	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Version
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// GetVersionInput is the input to the GetVersion function.
type GetVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the version number to fetch (required).
	ServiceVersion int
}

// GetVersion retrieves the specified resource.
func (c *Client) GetVersion(i *GetVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Version
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// UpdateVersionInput is the input to the UpdateVersion function.
type UpdateVersionInput struct {
	// Comment is a personal freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// UpdateVersion updates the specified resource.
func (c *Client) UpdateVersion(i *UpdateVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Version
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// ActivateVersionInput is the input to the ActivateVersion function.
type ActivateVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// Environment is the Fastly environment to activate this version to (optional).
	Environment string
}

// ActivateVersion activates the given version.
func (c *Client) ActivateVersion(i *ActivateVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	components := []string{"service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "activate"}
	if i.Environment != "" {
		components = append(components, i.Environment)
	}

	path := ToSafeURL(components...)

	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Version
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// DeactivateVersionInput is the input to the DeactivateVersion function.
type DeactivateVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// Environment is the Fastly environment to deactivate this version from (optional).
	Environment string
}

// DeactivateVersion deactivates the given version.
func (c *Client) DeactivateVersion(i *DeactivateVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	components := []string{"service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "deactivate"}
	if i.Environment != "" {
		components = append(components, i.Environment)
	}

	path := ToSafeURL(components...)

	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Version
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// CloneVersionInput is the input to the CloneVersion function.
type CloneVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// CloneVersion creates a clone of the specified version.
//
// Returns a new configuration version with all the same configuration options,
// but an incremented number.
func (c *Client) CloneVersion(i *CloneVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "clone")
	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Version
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// ValidateVersionInput is the input to the ValidateVersion function.
type ValidateVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ValidateVersion validates the specified resource.
func (c *Client) ValidateVersion(i *ValidateVersionInput) (bool, string, error) {
	var msg string

	if i.ServiceID == "" {
		return false, msg, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return false, msg, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "validate")
	resp, err := c.Get(path, nil)
	if err != nil {
		return false, msg, err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return false, msg, err
	}

	msg = r.Msg
	return r.Ok(), msg, nil
}

// LockVersionInput is the input to the LockVersion function.
type LockVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// LockVersion locks the specified version.
func (c *Client) LockVersion(i *LockVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "lock")
	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Version
	if err := DecodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}
