package suggest

// Suggestions is the API response structure for the suggest endpoint.
type Suggestions struct {
	Results []Suggestion `json:"results"`
}

// Suggestion represents an individual suggestion.
type Suggestion struct {
	// Domain is the full domain name suggestion
	Domain string `json:"domain"`
	// Subdomain is the portion of the domain before the zone
	Subdomain string `json:"subdomain"`
	// Zone is the top level domain or registered domain portion (e.g ".com")
	Zone string `json:"zone"`
}
