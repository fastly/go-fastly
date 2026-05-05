package dnszones

import "github.com/fastly/go-fastly/v14/fastly"

// Zone is the API response structure for the create, update and get operations.
type Zone struct {
	// ID is the zone identifier (UUID).
	ID *string `json:"id"`
	// Name is the domain name for your zone.
	Name *string `json:"name"`
	// Description is a freeform descriptive note.
	Description *string `json:"description"`
	// Type is the type of the zone.
	Type *string `json:"type"`
	// Serial is serial number of the zone's SOA record.
	Serial *string `json:"serial"`
	// Nameservers is an array of the nameserver hostnames assigned to the zone.
	Nameservers []string `json:"nameservers"`
	// XfrConfigInbound contains all attributes associated with inbound zone transfers.
	XfrConfigInbound *XfrConfigInbound `json:"xfr_config_inbound"`
	// CreatedAt is the date and time the zone was created.
	CreatedAt *string `json:"created_at"`
	// UpdatedAt is the date and time the zone was last updated.
	UpdatedAt *string `json:"updated_at"`
}

// XfrConfigInbound contains all attributes associated with inbound zone transfers.
type XfrConfigInbound struct {
	// Primaries is an array of the primary DNS server objects associated with inbound zone transfers.
	Primaries []Primary `json:"primaries"`
	// NotifyIPAddresses are IP addresses where Primary DNS servers can send NOTIFY messages.
	NotifyIPAddresses *NotifyIPAddresses `json:"notify_ip_addresses"`
	// InboundTSIGKeyID is the ID of the TSIG key used to secure inbound zone transfers.
	InboundTSIGKeyID *string `json:"inbound_tsig_key_id"`
}

// XfrConfigInboundInput specifies the inbound zone transfer attributes for create and update requests.
type XfrConfigInboundInput struct {
	// Primaries is an array of primary DNS server objects for inbound zone transfers.
	Primaries []Primary `json:"primaries,omitempty"`
	// InboundTSIGKeyID is the ID of the TSIG key used to secure inbound zone transfers.
	InboundTSIGKeyID *fastly.Nullable[string] `json:"inbound_tsig_key_id,omitempty"`
}

// Primary represents a primary DNS server for inbound zone transfers.
type Primary struct {
	// Address is an IPv4 address for the Primary DNS Server.
	Address *string `json:"address"`
	// Description is a description of the Primary DNS server.
	Description *string `json:"description"`
}

// NotifyIPAddresses contains IP addresses where primary DNS servers can send NOTIFY messages.
type NotifyIPAddresses struct {
	// IPv4 are IPv4 addresses where Primary DNS servers can send NOTIFY messages.
	IPv4 []string `json:"ipv4"`
}

// Zones is the paginated API response for listing zones.
type Zones struct {
	// Data is the list of zones.
	Data []Zone `json:"data"`
	// Meta contains pagination metadata.
	Meta MetaZones `json:"meta"`
}

// MetaZones is a subset of the Zone response structure.
type MetaZones struct {
	// NextCursor is the cursor value to use in the next request to retrieve the next page.
	NextCursor *string `json:"next_cursor"`
	// Limit is the maximum number of results returned.
	Limit *int `json:"limit"`
	// Sort is the order in which the results are listed.
	Sort *string `json:"sort"`
	// Total is the total number of zones.
	Total *int `json:"total"`
}
