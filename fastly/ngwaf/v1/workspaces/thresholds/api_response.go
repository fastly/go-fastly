package thresholds

import "time"

// Threshold is the API response structure for the create, update, and
// get operations.
type Threshold struct {
	// Action to take when threshold is exceeded. Must be one of
	// `block` or `log`.
	Action string `json:"action"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// DontNotify indicates whether to silence notifications when
	// action is taken. Defaults to false.
	DontNotify bool `json:"dont_notify"`
	// Duration is the duration the action is in place. Default
	// duration is 0.
	Duration int `json:"duration"`
	// Enabled is whether this threshold is active. Defaults to
	// false.
	Enabled bool `json:"enabled"`
	// Interval is the threshold interval in seconds. Must be one
	// of 60, 600, or 36000.
	Interval int `json:"interval"`
	// Limit is the threshold limit. Must be between than 1 and
	// less than 10000 inclusive.
	Limit int `json:"limit"`
	// Name is the threshold name.
	Name string `json:"name"`
	// Signal is the name of the signal this threshold is acting
	// on.
	Signal string `json:"signal"`
	// ThresholdID is the threshold identifier.
	ThresholdID string `json:"id"`
}

// Thresholds is the API response structure for the list Thresholds
// operation.
type Thresholds struct {
	// Data is the list of returned thresholds.
	Data []Threshold `json:"data"`
	// Meta is the information for total thresholds.
	Meta MetaThresholds `json:"meta"`
}

// MetaThresholds is a subset of the list thresholds response
// structure.
type MetaThresholds struct {
	// Limit is the limit of thresholds.
	Limit int `json:"limit"`
}
