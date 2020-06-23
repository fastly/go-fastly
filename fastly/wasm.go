package fastly

import (
	"fmt"
	"io"
	"time"
)

type WASMPackageMetadata struct {
	Name        string
	Description string
	Authors     []string
	Language    string
	Size        int
	HashSum     string
}

type WASMPackage struct {
	ID        string
	ServiceID string `mapstructure:"service_id"`
	Version   int
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
	Metadata  WASMPackageMetadata
}

// GetWASMPackageInput is used as input to the GetWASMPackage function.
type GetWASMPackageInput struct {
	// Service is the ID of the service.
	// Version is the specific configuration version.
	// Both fields are required.
	Service string `mapstructure:"service_id"`
	Version int    `mapstructure:"version"`
}

// GetWASMPackage retrieves WASM package information for the given service and version
func (c *Client) GetWASMPackage(i *GetWASMPackageInput) (*WASMPackage, error) {
	path, err := MakeWASMPackagePath(i.Service, i.Version)
	if err != nil {
		return nil, err
	}

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	return PopulateWASMPackage(resp.Body)
}

// UpdateWASMPackageInput is used as input to the UpdateWASMPackage function.
type UpdateWASMPackageInput struct {
	// Service is the ID of the service.
	// Version is the specific configuration version.
	// Both fields are required.
	Service string `mapstructure:"service_id"`
	Version int    `mapstructure:"version"`

	// PackagePath is the local filesystem path to the WASM package to upload
	PackagePath string
}

// UpdateWASMPackage updates a WASM package for a specific version.
func (c *Client) UpdateWASMPackage(i *UpdateWASMPackageInput) (*WASMPackage, error) {

	path, err := MakeWASMPackagePath(i.Service, i.Version)
	if err != nil {
		return nil, err
	}

	resp, err := c.PutFormFile(path, i, i.PackagePath, nil)
	if err != nil {
		return nil, err
	}

	return PopulateWASMPackage(resp.Body)
}

func MakeWASMPackagePath(Service string, Version int) (string, error) {
	if Service == "" {
		return "", ErrMissingService
	}
	if Version == 0 {
		return "", ErrMissingVersion
	}
	return fmt.Sprintf("/service/%s/version/%d/package", Service, Version), nil
}

func PopulateWASMPackage(body io.ReadCloser) (*WASMPackage, error) {
	var p *WASMPackage
	if err := decodeBodyMap(body, &p); err != nil {
		return nil, err
	}
	return p, nil
}
