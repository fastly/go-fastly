package status

// Status is the API response structure for the status endpoint.
type Status struct {
	// Domain is the domain whose status is being checked.
	Domain string `json:"domain,omitempty"`
	// Zone is the top level domain or registered domain portion (e.g ".com")
	Zone string `json:"zone,omitempty"`
	// Status is a space-delimited list of status types for the given domain
	// in increasing order of precedence.
	// The right-most value can be considered the most important value.
	Status string `json:"status,omitempty"`
	// Scope reflects the type of availability check that was performed.
	// Scope will only be present if an estimated status check was performed.
	// `estimate` provides DNS and aftermarket availability.
	Scope *Scope `json:"scope,omitempty"`
	// Tags is a space-delimited string containing the varying tags associated with the domain.
	Tags string `json:"tags,omitempty"`
	// Offers if present, contains a list of offers from domain aftermarket vendors.
	Offers []Offer `json:"offers,omitempty"`
}

// Offer represents an offer from an aftermarket vendor for a given domain.
type Offer struct {
	// Vendor is the name of the aftermarket vendor.
	Vendor string `json:"vendor,omitempty"`
	// Price is the price of the domain from the aftermarket vendor.
	Price string `json:"price,omitempty"`
	// Currency is the currency for the aftermarket offer.
	// A three-letter country currency code.
	Currency string `json:"currency,omitempty"`
}
