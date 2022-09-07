package fastly

import (
	"fmt"
	"net/url"
	"strings"
)

// Purge is a response from a purge request.
type Purge struct {
	// Status is the status of the purge, usually "ok".
	Status string `mapstructure:"status"`

	// ID is the unique ID of the purge request.
	ID string `mapstructure:"id"`
}

// PurgeInput is used as input to the Purge function.
type PurgeInput struct {
	// URL is the URL to purge (required).
	URL string

	// Soft performs a soft purge.
	Soft bool
}

// Purge instantly purges an individual URL.
func (c *Client) Purge(i *PurgeInput) (*Purge, error) {
	if i.URL == "" {
		return nil, ErrMissingURL
	}

	ro := &RequestOptions{
		Parallel: true,
	}
	if i.Soft {
		ro.Headers = map[string]string{
			"Fastly-Soft-Purge": "1",
		}
	}

	var err error
	ro.Params, err = constructRequestOptionsParam(i.URL)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post("purge/"+i.URL, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *Purge
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// constructRequestOptionsParam prevents Client.RawRequest from incorrectly
// clearing the URL query string, by manually constructing a ro.Params.
func constructRequestOptionsParam(us string) (map[string]string, error) {
	m := make(map[string]string)
	u, err := url.Parse(us)
	if err != nil {
		return nil, err
	}
	v := u.Query()
	// NOTE: we can't coerce a url.Values into the underlying map[string]string
	// type, so we have to manually loop over the url.Values and copy the
	// key/value pairs into a new map instance.
	for k, v := range v {
		m[k] = v[0]
	}
	return m, nil
}

// PurgeKeyInput is used as input to the PurgeKey function.
type PurgeKeyInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// Key is the key to purge (required).
	Key string

	// Soft performs a soft purge.
	Soft bool
}

// PurgeKey instantly purges a particular service of items tagged with a key.
func (c *Client) PurgeKey(i *PurgeKeyInput) (*Purge, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.Key == "" {
		return nil, ErrMissingKey
	}

	path := fmt.Sprintf("/service/%s/purge/%s", i.ServiceID, i.Key)

	ro := new(RequestOptions)
	ro.Parallel = true
	req, err := c.RawRequest("POST", path, ro)
	if err != nil {
		return nil, err
	}

	if i.Soft {
		req.Header.Set("Fastly-Soft-Purge", "1")
	}

	resp, err := checkResp(c.HTTPClient.Do(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *Purge
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// PurgeKeysInput is used as input to the PurgeKeys function.
type PurgeKeysInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// Keys are the keys to purge (required).
	Keys []string

	// Soft performs a soft purge.
	Soft bool
}

// PurgeKeys instantly purges a particular service of items tagged with a key.
func (c *Client) PurgeKeys(i *PurgeKeysInput) (map[string]string, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if len(i.Keys) == 0 {
		return nil, ErrMissingKeys
	}

	path := fmt.Sprintf("/service/%s/purge", i.ServiceID)

	ro := new(RequestOptions)
	ro.Parallel = true
	req, err := c.RawRequest("POST", path, ro)
	if err != nil {
		return nil, err
	}

	if i.Soft {
		req.Header.Set("Fastly-Soft-Purge", "1")
	}

	req.Header.Set("Surrogate-Key", strings.Join(i.Keys, " "))

	resp, err := checkResp(c.HTTPClient.Do(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r map[string]string
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// PurgeAllInput is used as input to the Purge function.
type PurgeAllInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// PurgeAll instantly purges everything from a service.
func (c *Client) PurgeAll(i *PurgeAllInput) (*Purge, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/purge_all", i.ServiceID)
	req, err := c.RawRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.HTTPClient.Do(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *Purge
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return nil, err
	}
	return r, nil
}
