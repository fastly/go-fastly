package computeacls

// ComputeACL is the API response structure for the create and describe operations.
type ComputeACL struct {
	// Name is an ACL name.
	Name string `json:"name"`
	// ComputeACLID is an ACL Identifier.
	ComputeACLID string `json:"id"`
}

// ComputeACLs is the API response structure for the list compute ACLs operation.
type ComputeACLs struct {
	// Data is the list of returned compute ACLs.
	Data []ComputeACL `json:"data"`
	// Meta is the information for total compute ACLs.
	Meta MetaACLs `json:"meta"`
}

// MetaACLs is a subset of the ComputeACLs response structure.
type MetaACLs struct {
	// Total is the sum of compute ACLs.
	Total int `json:"total"`
}

// ComputeACLEntry is the API response structure for the lookup operation.
type ComputeACLEntry struct {
	// Prefix is an IP prefix defined in Classless Inter-Domain Routing (CIDR) format.
	Prefix string `json:"prefix"`
	// Action is one of "ALLOW" or "BLOCK".
	Action string `json:"action"`
}

// ComputeACLEntries is the API response structure for the list compute ACL entries operation.
type ComputeACLEntries struct {
	// Entries is the list of returned compute ACL entries.
	Entries []ComputeACLEntry
	// Meta is the information for pagination.
	Meta MetaEntries `json:"meta"`
}

// MetaEntries is a subset of the ComputeACLs response structure.
type MetaEntries struct {
	// Limit is how many results are included in this response.
	Limit int `json:"limit"`
	// NextCursor is the cursor value used to retrieve the next page.
	NextCursor string `json:"next_cursor"`
}
