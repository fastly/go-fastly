package fastly

import "context"

// IPAddrs is a list of IP addresses returned by the Fastly API.
type IPAddrs []string

// AllIPs returns the lists of public IPv4 and IPv6 addresses for Fastly's network.
func (c *Client) AllIPs(ctx context.Context) (v4, v6 IPAddrs, err error) {
	resp, err := c.Get(ctx, "/public-ip-list", CreateRequestOptions())
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var m map[string][]string
	if err := DecodeBodyMap(resp.Body, &m); err != nil {
		return nil, nil, err
	}

	return m["addresses"], m["ipv6_addresses"], nil
}

// IPs returns the list of public IPv4 addresses for Fastly's network.
func (c *Client) IPs(ctx context.Context) (IPAddrs, error) {
	v4, _, err := c.AllIPs(ctx)
	if err != nil {
		return nil, err
	}

	return v4, nil
}

// IPsV6 returns the list of public IPv6 addresses for Fastly's network.
func (c *Client) IPsV6(ctx context.Context) (IPAddrs, error) {
	_, v6, err := c.AllIPs(ctx)
	if err != nil {
		return nil, err
	}

	return v6, nil
}
