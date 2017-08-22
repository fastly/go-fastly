package fastly

import "testing"

// testWafID is the ID of one of the WAF firewall objects.
var testWafID string

func TestClient_GetFirewallObjects(t *testing.T) {
	var err error
	var tv *Version
	record(t, "waf/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	var firewallObjects *FirewallObjects
	record(t, "waf/get", func(c *Client) {
		firewallObjects, err = c.GetFirewallObjects(&GetFirewallObjectsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(firewallObjects.Data) < 1 {
		t.Errorf("did not properly list firewall objects: %v", firewallObjects)
	}

	testServiceID = firewallObjects.Data[0].ID
}

func TestClient_GetFirewallObject(t *testing.T) {
	var err error
	var tv *Version
	record(t, "waf/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	var firewallObject *FirewallObject
	record(t, "waf/getObject", func(c *Client) {
		firewallObject, err = c.GetFirewallObject(&GetFirewallObjectInput{
			Service: testServiceID,
			Version: tv.Number,
			WafID:   testWafID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if firewallObject.Data.ID != testWafID {
		t.Errorf("did not properly get firewall object: %v", firewallObject)
	}
}
