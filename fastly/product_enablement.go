package fastly

import (
	"fmt"
)

// ProductEnablement represents a response from the Fastly API.
type ProductEnablement struct {
	Product ProductEnablementNested `mapstructure:"product"`
	Service ProductEnablementNested `mapstructure:"service"`
}

type ProductEnablementNested struct {
	ID     string `mapstructure:"id,omitempty"`
	Object string `mapstructure:"object,omitempty"`
}

// Product is a base for the different product variants.
type Product int64

func (p Product) String() string {
	switch p {
	case ProductBrotliCompression:
		return "brotli_compression"
	case ProductDomainInspector:
		return "domain_inspector"
	case ProductFanout:
		return "fanout"
	case ProductImageOptimizer:
		return "image_optimizer"
	case ProductOriginInspector:
		return "origin_inspector"
	case ProductWebSockets:
		return "websockets"
	}
	return "unknown"
}

const (
	ProductUndefined Product = iota
	ProductBrotliCompression
	ProductDomainInspector
	ProductFanout
	ProductImageOptimizer
	ProductOriginInspector
	ProductWebSockets
)

// ProductEnablementInput is used as input to the various product API functions.
type ProductEnablementInput struct {
	// ProductID is the ID of the product and is constrained by the Product type (required).
	ProductID Product
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// GetProduct retrieves the details of the product enabled on the service.
func (c *Client) GetProduct(i *ProductEnablementInput) (*ProductEnablement, error) {
	if i.ProductID == ProductUndefined {
		return nil, ErrMissingProductID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/enabled-products/%s/services/%s", i.ProductID, i.ServiceID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *ProductEnablement
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}

	return h, nil
}

// EnableProduct enables the specified product on the service.
func (c *Client) EnableProduct(i *ProductEnablementInput) (*ProductEnablement, error) {
	if i.ProductID == ProductUndefined {
		return nil, ErrMissingProductID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/enabled-products/%s/services/%s", i.ProductID, i.ServiceID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var http3 *ProductEnablement
	if err := decodeBodyMap(resp.Body, &http3); err != nil {
		return nil, err
	}
	return http3, nil
}

// DisableProduct disables the specified product on the service.
func (c *Client) DisableProduct(i *ProductEnablementInput) error {
	if i.ProductID == ProductUndefined {
		return ErrMissingProductID
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := fmt.Sprintf("/enabled-products/%s/services/%s", i.ProductID, i.ServiceID)
	_, err := c.Delete(path, nil)
	return err
}
