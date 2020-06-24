package fastly

import (
	"fmt"
	"io"
	"time"
)

type WasmPackageMetadata struct {
	Name        string
	Description string
	Authors     []string
	Language    string
	Size        int
	HashSum     string
}

type WasmPackage struct {
	ID        string
	ServiceID string `mapstructure:"service_id"`
	Version   int
	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
	DeletedAt *time.Time `mapstructure:"deleted_at"`
	Metadata  WasmPackageMetadata
}

// GetWasmPackageInput is used as input to the GetWasmPackage function.
type GetWasmPackageInput struct {
	// Service is the ID of the service.
	// Version is the specific configuration version.
	// Both fields are required.
	Service string `mapstructure:"service_id"`
	Version int    `mapstructure:"version"`
}

// GetWasmPackage retrieves Wasm package information for the given service and version
func (c *Client) GetWasmPackage(i *GetWasmPackageInput) (*WasmPackage, error) {
	path, err := MakeWasmPackagePath(i.Service, i.Version)
	if err != nil {
		return nil, err
	}

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	return PopulateWasmPackage(resp.Body)
}

// UpdateWasmPackageInput is used as input to the UpdateWasmPackage function.
type UpdateWasmPackageInput struct {
	// Service is the ID of the service.
	// Version is the specific configuration version.
	// Both fields are required.
	Service string `mapstructure:"service_id"`
	Version int    `mapstructure:"version"`

	// PackagePath is the local filesystem path to the Wasm package to upload
	PackagePath string
}

// UpdateWasmPackage updates a Wasm package for a specific version.
func (c *Client) UpdateWasmPackage(i *UpdateWasmPackageInput) (*WasmPackage, error) {

	path, err := MakeWasmPackagePath(i.Service, i.Version)
	if err != nil {
		return nil, err
	}

	resp, err := c.PutFormFile(path, i, i.PackagePath, nil)
	if err != nil {
		return nil, err
	}

	return PopulateWasmPackage(resp.Body)
}

func MakeWasmPackagePath(Service string, Version int) (string, error) {
	if Service == "" {
		return "", ErrMissingService
	}
	if Version == 0 {
		return "", ErrMissingVersion
	}
	return fmt.Sprintf("/service/%s/version/%d/package", Service, Version), nil
}

func PopulateWasmPackage(body io.ReadCloser) (*WasmPackage, error) {
	var p *WasmPackage
	if err := decodeBodyMap(body, &p); err != nil {
		return nil, err
	}
	return p, nil
}
