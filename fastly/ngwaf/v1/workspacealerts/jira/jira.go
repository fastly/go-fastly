package jira

// Integration is the config object for integration type jira.
type Integration struct {
	// Host is the name of the Jira instnace (required).
	Host *string `json:"host,omitempty"`
	// IssueType is the Jira issue type associated with the ticket (optional).
	IssueType *string `json:"issue_type,omitempty"`
	// Key is the Jira API key / secret field (required).
	Key *string `json:"key,omitempty"`
	// Project specifies the Jira project where the issue will be created (required).
	Project *string `json:"project,omitempty"`
	// Username is the Jira username of the user who created the ticket (required).
	Username *string `json:"username,omitempty"`
}