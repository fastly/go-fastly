package fastly

import (
	"fmt"
	"io"
	"time"
)

// Package is a container for data returned about a package.
type Package struct {
	ID             string
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`
	Metadata       PackageMetadata
	CreatedAt      *time.Time `mapstructure:"created_at"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
}

// Package is a container for metadata returned about a package.
// It is a separate struct to allow correct serialisation by mapstructure -
// the raw data is returned as a json sub-block.
type PackageMetadata struct {
	Name        string
	Description string
	Authors     []string
	Language    string
	Size        int64
	HashSum     string
}

// GetPackageInput is used as input to the GetPackage function.
type GetPackageInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string `mapstructure:"service_id"`

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `mapstructure:"version"`
}

// GetPackage retrieves  package information for the given service and version.
func (c *Client) GetPackage(i *GetPackageInput) (*Package, error) {
	path, err := MakePackagePath(i.ServiceID, i.ServiceVersion)
	if err != nil {
		return nil, err
	}

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	return PopulatePackage(resp.Body)
}

// UpdatePackageInput is used as input to the UpdatePackage function.
type UpdatePackageInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string `mapstructure:"service_id"`

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `mapstructure:"version"`

	// PackagePath is the local filesystem path to the package to upload.
	PackagePath string
}

// UpdatePackage updates a package for a specific version.
func (c *Client) UpdatePackage(i *UpdatePackageInput) (*Package, error) {

	urlPath, err := MakePackagePath(i.ServiceID, i.ServiceVersion)
	if err != nil {
		return nil, err
	}

	resp, err := c.PutFormFile(urlPath, i.PackagePath, "package", nil)
	if err != nil {
		return nil, err
	}

	return PopulatePackage(resp.Body)
}

// MakePackagePath ensures we create the correct REST path for referencing packages in the API.
func MakePackagePath(ServiceID string, ServiceVersion int) (string, error) {
	if ServiceID == "" {
		return "", ErrMissingServiceID
	}
	if ServiceVersion == 0 {
		return "", ErrMissingServiceVersion
	}
	return fmt.Sprintf("/service/%s/version/%d/package", ServiceID, ServiceVersion), nil
}

// PopulatePackage encapsulates the decoding of returned package data.
func PopulatePackage(body io.ReadCloser) (*Package, error) {
	var p *Package
	if err := decodeBodyMap(body, &p); err != nil {
		return nil, err
	}
	return p, nil
}
