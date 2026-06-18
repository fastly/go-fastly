package tsigkeys

// TSIGKey is the API response structure for the create, update and get operations.
type TSIGKey struct {
	// ID is the TSIG Key Identifier (UUID).
	ID *string `json:"id"`
	// Name is the name of the TSIG key.
	Name *string `json:"name"`
	// Description is a freeform descriptive note.
	Description *string `json:"description"`
	// Algorithm is the algorithm of the TSIG key.
	Algorithm *string `json:"algorithm"`
	// CreatedAt is the date and time the TSIG key was created.
	CreatedAt *string `json:"created_at"`
	// UpdatedAt is the date and time the TSIG key was last updated.
	UpdatedAt *string `json:"updated_at"`
}

// TSIGKeys is the paginated API response for listing TSIG keys.
type TSIGKeys struct {
	// Data is the list of TSIG keys.
	Data []TSIGKey `json:"data"`
	// Meta contains pagination metadata.
	Meta MetaTSIGKeys `json:"meta"`
}

// MetaTSIGKeys is a subset of the TSIGKey response structure.
type MetaTSIGKeys struct {
	// NextCursor is the cursor value to use in the next request to retrieve the next page.
	NextCursor *string `json:"next_cursor"`
	// Limit is the maximum number of results returned.
	Limit *int `json:"limit"`
	// Sort is the order in which the results are listed.
	Sort *string `json:"sort"`
	// Total is the total number of TSIG keys.
	Total *int `json:"total"`
}
