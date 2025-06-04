package signals

// Reason is the signal that corresponds to the reason an event was triggered.
type Reason struct {
	// Signal ID is the ID of the signal that triggered the event
	SignalID string `json:"signal_id"`
	// Count is the number of times this signal was detected
	Count int `json:"count"`
}
