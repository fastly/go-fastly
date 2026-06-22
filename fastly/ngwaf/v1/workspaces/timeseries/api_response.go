package timeseries

// TimeSeries is the API response structure for the get time series operations.
type TimeSeries struct {
	// Data is the list of returned time series data points.
	Data []map[string]any `json:"data"`
	// Meta is the information for total time series.
	Meta MetaTimeSeries `json:"meta"`
}

// MetaTimeSeries is a subset of the TimeSeries response structure.
type MetaTimeSeries struct {
	// Total is the sum of TimeSeries.
	Total int `json:"total"`
}
