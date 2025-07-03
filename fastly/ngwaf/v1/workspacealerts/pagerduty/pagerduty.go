package pagerduty

// Integration is the config object for integration type pagerduty.
type Integration struct {
	// Key is the PagerDuty integration key (required).
	Key *string `json:"key,omitempty"`
}