package fastly

import (
	"fmt"
	"sort"
	"time"
)

// Version represents a distinct configuration version.
type Version struct {
	Number    int        `mapstructure:"number"`
	Comment   string     `mapstructure:"comment"`
	ServiceID string     `mapstructure:"service_id"`
	Active    bool       `mapstructure:"active"`
	Locked    bool       `mapstructure:"locked"`
	Deployed  bool       `mapstructure:"deployed"`
	Staging   bool       `mapstructure:"staging"`
	Testing   bool       `mapstructure:"testing"`
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
}

// versionsByNumber is a sortable list of versions. This is used by the version
// `List()` function to sort the API responses.
type versionsByNumber []*Version

// Len, Swap, and Less implement the sortable interface.
func (s versionsByNumber) Len() int      { return len(s) }
func (s versionsByNumber) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s versionsByNumber) Less(i, j int) bool {
	return s[i].Number < s[j].Number
}

// ListVersionsInput is the input to the ListVersions function.
type ListVersionsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// ListVersions returns the full list of all versions of the given service.
func (c *Client) ListVersions(i *ListVersionsInput) ([]*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/version", i.ServiceID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var e []*Version
	if err := decodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	sort.Sort(versionsByNumber(e))

	return e, nil
}

// LatestVersionInput is the input to the LatestVersion function.
type LatestVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// LatestVersion fetches the latest version. If there are no versions, this
// function will return nil (but not an error).
func (c *Client) LatestVersion(i *LatestVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	list, err := c.ListVersions(&ListVersionsInput{ServiceID: i.ServiceID})
	if err != nil {
		return nil, err
	}
	if len(list) < 1 {
		return nil, nil
	}

	e := list[len(list)-1]
	return e, nil
}

// CreateVersionInput is the input to the CreateVersion function.
type CreateVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// A personal freeform descriptive note.
	Comment string `form:"comment,omitempty"`
}

// CreateVersion constructs a new version. Note that `CloneVersion` is
// preferred in almost all scenarios, since `Create()` creates a _blank_
// configuration where `Clone()` builds off of an existing configuration.
func (c *Client) CreateVersion(i *CreateVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/version", i.ServiceID)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var e *Version
	if err := decodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// GetVersionInput is the input to the GetVersion function.
type GetVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// SrrviceVersion is the version number to fetch (required).
	ServiceVersion int
}

// GetVersion fetches a version with the given information.
func (c *Client) GetVersion(i *GetVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var e *Version
	if err := decodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// UpdateVersionInput is the input to the UpdateVersion function.
type UpdateVersionInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// A personal freeform descriptive note.
	Comment *string `form:"comment,omitempty"`
}

// UpdateVersion updates the given version
func (c *Client) UpdateVersion(i *UpdateVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d", i.ServiceID, i.ServiceVersion)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var e *Version
	if err := decodeBodyMap(resp.Body, &e); err != nil {
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
}

// ActivateVersion activates the given version.
func (c *Client) ActivateVersion(i *ActivateVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/activate", i.ServiceID, i.ServiceVersion)
	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}

	var e *Version
	if err := decodeBodyMap(resp.Body, &e); err != nil {
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
}

// DeactivateVersion deactivates the given version.
func (c *Client) DeactivateVersion(i *DeactivateVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/deactivate", i.ServiceID, i.ServiceVersion)
	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}

	var e *Version
	if err := decodeBodyMap(resp.Body, &e); err != nil {
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

// CloneVersion creates a clone of the version with and returns a new
// configuration version with all the same configuration options, but an
// incremented number.
func (c *Client) CloneVersion(i *CloneVersionInput) (*Version, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/clone", i.ServiceID, i.ServiceVersion)
	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}

	var e *Version
	if err := decodeBodyMap(resp.Body, &e); err != nil {
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

// ValidateVersion validates if the given version is okay.
func (c *Client) ValidateVersion(i *ValidateVersionInput) (bool, string, error) {
	var msg string

	if i.ServiceID == "" {
		return false, msg, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return false, msg, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/validate", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return false, msg, err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
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

	path := fmt.Sprintf("/service/%s/version/%d/lock", i.ServiceID, i.ServiceVersion)
	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}

	var e *Version
	if err := decodeBodyMap(resp.Body, &e); err != nil {
		return nil, err
	}
	return e, nil
}
