package fastly

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
)

// Package is a container for data returned about a package.
type Package struct {
	CreatedAt      *time.Time       `mapstructure:"created_at"`
	DeletedAt      *time.Time       `mapstructure:"deleted_at"`
	Metadata       *PackageMetadata `mapstructure:"metadata"`
	PackageID      *string          `mapstructure:"id"`
	ServiceID      *string          `mapstructure:"service_id"`
	ServiceVersion *int             `mapstructure:"version"`
	UpdatedAt      *time.Time       `mapstructure:"updated_at"`
}

// PackageMetadata is a container for metadata returned about a package.
// It is a separate struct to allow correct serialisation by mapstructure -
// the raw data is returned as a json sub-block.
type PackageMetadata struct {
	Authors     []string `mapstructure:"authors"`
	Description *string  `mapstructure:"description"`
	FilesHash   *string  `mapstructure:"files_hash"`
	HashSum     *string  `mapstructure:"hashsum"`
	Language    *string  `mapstructure:"language"`
	Name        *string  `mapstructure:"name"`
	Size        *int64   `mapstructure:"size"`
}

// GetPackageInput is used as input to the GetPackage function.
type GetPackageInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string `mapstructure:"service_id"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `mapstructure:"version"`
}

// GetPackage retrieves the specified resource.
func (c *Client) GetPackage(i *GetPackageInput) (*Package, error) {
	path, err := MakePackagePath(i.ServiceID, i.ServiceVersion)
	if err != nil {
		return nil, err
	}

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return PopulatePackage(resp.Body)
}

// UpdatePackageInput is used as input to the UpdatePackage function.
type UpdatePackageInput struct {
	// PackagePath is the local filesystem path to the package to upload.
	PackagePath *string
	// PackageContent is the data in raw of the package to upload.
	PackageContent []byte
	// ServiceID is the ID of the service (required).
	ServiceID string `mapstructure:"service_id"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `mapstructure:"version"`
}

// UpdatePackage updates the specified resource.
func (c *Client) UpdatePackage(i *UpdatePackageInput) (*Package, error) {
	urlPath, err := MakePackagePath(i.ServiceID, i.ServiceVersion)
	if err != nil {
		return nil, err
	}

	var body io.ReadCloser
	switch {
	case i.PackagePath != nil:
		resp, err := c.PutFormFile(urlPath, *i.PackagePath, "package", nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body = resp.Body
	case len(i.PackageContent) != 0:
		resp, err := c.PutFormFileFromReader(urlPath, "package.tar.gz", bytes.NewReader(i.PackageContent), "package", nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body = resp.Body
	default:
		return nil, errors.New("missing package file path or content")
	}

	return PopulatePackage(body)
}

// MakePackagePath ensures we create the correct REST path for referencing packages in the API.
func MakePackagePath(serviceID string, serviceVersion int) (string, error) {
	if serviceID == "" {
		return "", ErrMissingServiceID
	}
	if serviceVersion == 0 {
		return "", ErrMissingServiceVersion
	}
	return fmt.Sprintf("/service/%s/version/%d/package", serviceID, serviceVersion), nil
}

// PopulatePackage encapsulates the decoding of returned package data.
func PopulatePackage(body io.ReadCloser) (*Package, error) {
	var p *Package
	if err := decodeBodyMap(body, &p); err != nil {
		return nil, err
	}
	return p, nil
}
