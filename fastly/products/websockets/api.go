// Code generated by 'service_linked_product' generator, DO NOT EDIT.

package websockets

import fastly "github.com/fastly/go-fastly/v9/fastly"

// Get gets the status of the WebSockets product on the service.
func Get(c *fastly.Client, serviceID string) (*fastly.ProductEnablement, error) {
	if serviceID == "" {
		return nil, fastly.ErrMissingServiceID
	}

	path := fastly.ToSafeURL("enabled-products", "v1", "websockets", "services", serviceID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *fastly.ProductEnablement
	if err := fastly.DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// Enable enables the WebSockets product on the service.
func Enable(c *fastly.Client, serviceID string) (*fastly.ProductEnablement, error) {
	if serviceID == "" {
		return nil, fastly.ErrMissingServiceID
	}

	path := fastly.ToSafeURL("enabled-products", "v1", "websockets", "services", serviceID)

	resp, err := c.Put(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *fastly.ProductEnablement
	if err := fastly.DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// Disable disables the WebSockets product on the service.
func Disable(c *fastly.Client, serviceID string) error {
	if serviceID == "" {
		return fastly.ErrMissingServiceID
	}

	path := fastly.ToSafeURL("enabled-products", "v1", "websockets", "services", serviceID)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
