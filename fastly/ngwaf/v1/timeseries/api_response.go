package timeseries

// TimeseriesDataPoint is a single data point in the timeseries.
type TimeseriesDataPoint struct {
	// Dimensions are the dimension values for this data point.
	Dimensions TimeseriesDimensions `json:"dimensions"`
	// Values are the metric values for this data point.
	Values []map[string]any `json:"values"`
}

// TimeseriesDimensions holds the grouping dimensions for a data point.
type TimeseriesDimensions struct {
	// Time is the timestamp for this data point.
	Time string `json:"time"`
	// Workspace is the workspace identifier.
	Workspace string `json:"workspace"`
}

// Timeseries is the API response structure for the list timeseries operation.
type Timeseries struct {
	// Data is the list of returned timeseries data points.
	Data []TimeseriesDataPoint `json:"data"`
	// Meta is the information for total timeseries.
	Meta MetaTimeseries `json:"meta"`
}

// MetaTimeseries is a subset of the Timeseries response structure.
type MetaTimeseries struct {
	// Total is the count of data points returned.
	Total int `json:"total"`
}
